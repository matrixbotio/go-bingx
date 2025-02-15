package bingxgo

import (
	"strconv"
)

type MarketClient struct {
	client *Client
}

func (c *MarketClient) GetKlines(symbol string, interval string, limit int) ([]Kline, error) {
	params := map[string]interface{}{
		"symbol":   symbol,
		"interval": interval,
		"limit":    strconv.Itoa(limit),
	}

	resp, err := c.client.sendRequest("GET", "/openApi/swap/v3/quote/klines", params)
	if err != nil {
		return nil, err
	}

	var klines []Kline
	err = json.Unmarshal(resp, &klines)
	return klines, err
}
