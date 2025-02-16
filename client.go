package bingxgo

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const baseURL = "https://open-api.bingx.com"

type Client struct {
	HTTPClient  *http.Client
	rateLimiter *RateLimiter

	APIKey          string
	SecretKey       string
	BaseURL         string
	Debug           bool
	BrokerSourceKey string
}

func NewClient(apiKey, secretKey string) *Client {
	return &Client{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		Debug:      false,
	}
}

func (c *Client) SetRateLimiter(rateLimiter *RateLimiter) *Client {
	c.rateLimiter = rateLimiter
	return c
}

func (c *Client) SetBrokerSourceKey(key string) *Client {
	c.BrokerSourceKey = key
	return c
}

func (c *Client) sendRequest(
	method string,
	endpoint string,
	params map[string]interface{},
	resultPointer any,
) error {
	if len(params) == 0 {
		return fmt.Errorf("params map is nil or empty")
	}
	if c.rateLimiter != nil {
		c.rateLimiter.Wait(endpoint)
	}

	// Build query parameters
	encodedParams, rawParams := c.buildParams(params)

	// Generate signature
	signature := c.generateSignature(rawParams)

	// Create full URL
	fullURL := c.buildURL(endpoint, encodedParams, signature)

	// Execute request
	body, err := c.executeRequest(method, fullURL)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, resultPointer); err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	return nil
}

func (c *Client) buildParams(params map[string]interface{}) (encoded, raw string) {
	var encodedBuilder, rawBuilder strings.Builder

	// Sort keys for consistent ordering
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build both encoded and raw params
	for _, k := range keys {
		value := fmt.Sprintf("%v", params[k])
		encodedValue := url.QueryEscape(value)
		encodedValue = strings.ReplaceAll(encodedValue, "+", "%20")

		encodedBuilder.WriteString(k + "=" + encodedValue + "&")
		rawBuilder.WriteString(k + "=" + value + "&")
	}

	// Add timestamp to both
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	encodedBuilder.WriteString("timestamp=" + timestamp)
	rawBuilder.WriteString("timestamp=" + timestamp)

	if c.Debug {
		log.Printf("Raw params: %s", rawBuilder.String())
	}

	return encodedBuilder.String(), rawBuilder.String()
}

func (c *Client) buildURL(endpoint, params, signature string) string {
	fullURL := fmt.Sprintf("%s%s?%s&signature=%s",
		c.BaseURL,
		endpoint,
		params,
		signature,
	)

	if c.Debug {
		log.Printf("Full URL: %s", fullURL)
	}

	return fullURL
}

func (c *Client) executeRequest(method, url string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("X-BX-APIKEY", c.APIKey)
	if c.BrokerSourceKey != "" {
		req.Header.Set("X-SOURCE-KEY", c.BrokerSourceKey)
	}

	if c.Debug {
		log.Printf("Request Headers: %v", req.Header)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if c.Debug {
		log.Printf("Response Body: %s", string(body))
	}

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp.StatusCode, body)
	}

	return body, nil
}

func (c *Client) handleErrorResponse(statusCode int, body []byte) ([]byte, error) {
	var apiErr APIError
	if err := json.Unmarshal(body, &apiErr); err != nil {
		return nil, fmt.Errorf("http status %d (%s), body: %s",
			statusCode,
			http.StatusText(statusCode),
			string(body),
		)
	}
	return nil, apiErr
}

func (c *Client) generateSignature(queryString string) string {
	h := hmac.New(sha256.New, []byte(c.SecretKey))
	h.Write([]byte(queryString))
	return hex.EncodeToString(h.Sum(nil))
}

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error, code: %d, message: %s", e.Code, e.Message)
}
