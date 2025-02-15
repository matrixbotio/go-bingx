package bingxgo

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type SpotClient struct {
	client *Client
}

func NewSpotClient(client *Client) SpotClient {
	return SpotClient{client: client}
}

func (c *SpotClient) GetBalance() ([]SpotBalance, error) {
	endpoint := "/openApi/spot/v1/account/balance"
	params := map[string]interface{}{
		"timestamp": time.Now().UnixMilli(),
	}

	resp, err := c.client.sendRequest("GET", endpoint, params)
	if err != nil {
		return nil, err
	}

	var bingXResponse BingXResponse[map[string][]SpotBalance]
	err = json.Unmarshal(resp, &bingXResponse)
	if err != nil {
		return nil, err
	}
	if err := bingXResponse.Error(); err != nil {
		return nil, err
	}
	return bingXResponse.Data["balances"], err
}

func (c *SpotClient) CreateOrder(order SpotOrderRequest) (*SpotOrderResponse, error) {
	endpoint := "/openApi/spot/v1/trade/order"
	params := map[string]interface{}{
		"symbol":   order.Symbol,
		"side":     string(order.Side),
		"type":     string(order.Type),
		"quantity": decimal.NewFromFloat(order.Quantity).String(),
		"price":    decimal.NewFromFloat(order.Price).String(),
	}
	if order.ClientOrderID != "" {
		params["newClientOrderId"] = order.ClientOrderID
	}

	resp, err := c.client.sendRequest("POST", endpoint, params)
	if err != nil {
		return nil, err
	}

	var bingXResponse BingXResponse[SpotOrderResponse]
	err = json.Unmarshal(resp, &bingXResponse)
	if err != nil {
		return nil, err
	}
	if err := bingXResponse.Error(); err != nil {
		return nil, err
	}
	return &bingXResponse.Data, err
}

func (c *SpotClient) CreateBatchOrders(orders []SpotOrderRequest, isSync bool) ([]SpotOrderResponse, error) {
	endpoint := "/openApi/spot/v1/trade/batchOrders"

	ordersJSON, err := json.Marshal(orders)
	if err != nil {
		return nil, err
	}
	params := map[string]interface{}{
		"data": string(ordersJSON),
		"sync": isSync,
	}

	resp, err := c.client.sendRequest("POST", endpoint, params)
	if err != nil {
		return nil, err
	}

	var bingXResponse BingXResponse[map[string][]SpotOrderResponse]
	err = json.Unmarshal(resp, &bingXResponse)
	if err != nil {
		return nil, err
	}
	if err := bingXResponse.Error(); err != nil {
		return nil, err
	}
	return bingXResponse.Data["orders"], err
}

func (c *SpotClient) GetOpenOrders(symbol string) ([]SpotOrder, error) {
	endpoint := "/openApi/spot/v1/trade/openOrders"
	params := map[string]interface{}{
		"symbol": symbol,
	}

	resp, err := c.client.sendRequest("GET", endpoint, params)
	if err != nil {
		return nil, err
	}

	var bingXResponse BingXResponse[map[string][]SpotOrder]
	err = json.Unmarshal(resp, &bingXResponse)
	if err != nil {
		return nil, err
	}
	if err := bingXResponse.Error(); err != nil {
		return nil, err
	}
	return bingXResponse.Data["orders"], err
}

func (c *SpotClient) CancelOrder(symbol string, orderId string) error {
	endpoint := "/openApi/spot/v1/trade/cancel"
	params := map[string]interface{}{
		"symbol":  symbol,
		"orderId": orderId,
	}

	resp, err := c.client.sendRequest("POST", endpoint, params)
	if err != nil {
		return err
	}
	var bingXResponse BingXResponse[any]
	err = json.Unmarshal(resp, &bingXResponse)
	if err != nil {
		return err
	}
	if err := bingXResponse.Error(); err != nil {
		return err
	}
	return nil
}

func (c *SpotClient) CancelOrderByClientOrderID(
	symbol string,
	clientOrderID string,
) error {
	endpoint := "/openApi/spot/v1/trade/cancel"
	params := map[string]interface{}{
		"symbol":        symbol,
		"clientOrderID": clientOrderID,
	}

	resp, err := c.client.sendRequest("POST", endpoint, params)
	if err != nil {
		return err
	}
	var bingXResponse BingXResponse[any]
	err = json.Unmarshal(resp, &bingXResponse)
	if err != nil {
		return err
	}
	if err := bingXResponse.Error(); err != nil {
		return err
	}
	return nil
}

func (c *SpotClient) CancelAllOpenOrders(symbol string) error {
	endpoint := "/openApi/spot/v1/trade/cancelOpenOrders"
	params := map[string]interface{}{
		"symbol": symbol,
	}

	resp, err := c.client.sendRequest("POST", endpoint, params)
	if err != nil {
		return err
	}
	var bingXResponse BingXResponse[any]
	err = json.Unmarshal(resp, &bingXResponse)
	if err != nil {
		return err
	}
	if err := bingXResponse.Error(); err != nil {
		return err
	}
	return err
}

func (c *SpotClient) GetOrder(symbol string, orderID int64) (*SpotOrder, error) {
	return c.getOrderData(map[string]interface{}{
		"symbol":    symbol,
		"orderId":   orderID,
		"timestamp": time.Now().UnixMilli(),
	})
}

func (c *SpotClient) GetOrderByClientOrderID(
	symbol string,
	clientOrderID string,
) (*SpotOrder, error) {
	return c.getOrderData(map[string]interface{}{
		"symbol":        symbol,
		"clientOrderID": clientOrderID,
		"timestamp":     time.Now().UnixMilli(),
	})
}

func (c *SpotClient) getOrderData(
	params map[string]interface{},
) (*SpotOrder, error) {
	endpoint := "/openApi/spot/v1/trade/order"

	resp, err := c.client.sendRequest("GET", endpoint, params)
	if err != nil {
		return nil, err
	}

	var bingXResponse BingXResponse[SpotOrder]
	err = json.Unmarshal(resp, &bingXResponse)
	if err != nil {
		return nil, err
	}
	if err := bingXResponse.Error(); err != nil {
		return nil, err
	}
	return &bingXResponse.Data, err
}

func (c *SpotClient) HistoryOrders(symbol string) ([]SpotOrder, error) {
	endpoint := "/openApi/spot/v1/trade/historyOrders"
	params := map[string]interface{}{
		"symbol": symbol,
	}

	resp, err := c.client.sendRequest("GET", endpoint, params)
	if err != nil {
		return nil, err
	}

	var bingXResponse BingXResponse[map[string][]SpotOrder]
	err = json.Unmarshal(resp, &bingXResponse)
	if err != nil {
		return nil, err
	}
	if err := bingXResponse.Error(); err != nil {
		return nil, err
	}
	return bingXResponse.Data["orders"], err
}

func (c *SpotClient) OrderBook(symbol string, limit int) (*OrderBook, error) {
	endpoint := "/openApi/spot/v1/market/depth"
	params := map[string]interface{}{
		"symbol": symbol,
	}
	if limit > 0 {
		params["limit"] = limit
	}

	resp, err := c.client.sendRequest("GET", endpoint, params)
	if err != nil {
		return nil, err
	}

	var bingXResponse BingXResponse[OrderBook]
	err = json.Unmarshal(resp, &bingXResponse)
	if err != nil {
		return nil, err
	}
	if err := bingXResponse.Error(); err != nil {
		return nil, err
	}
	return &bingXResponse.Data, err
}

func (c *SpotClient) GetSymbols(symbol ...string) ([]SymbolInfo, error) {
	endpoint := "/openApi/spot/v1/common/symbols"
	params := map[string]interface{}{
		"timestamp": time.Now().UnixMilli(),
	}
	if len(symbol) > 0 {
		params["symbol"] = symbol[0]
	}

	resp, err := c.client.sendRequest("GET", endpoint, params)
	if err != nil {
		return nil, err
	}

	var bingXResponse BingXResponse[SymbolInfos]
	err = json.Unmarshal(resp, &bingXResponse)
	if err != nil {
		return nil, err
	}
	if err := bingXResponse.Error(); err != nil {
		return nil, err
	}
	return bingXResponse.Data.Symbols, nil
}

func (c *SpotClient) GetHistoricalKlines(
	symbol string,
	interval string,
	limit int64,
) ([]KlineData, error) {
	endpoint := "/openApi/market/his/v1/kline"
	params := map[string]interface{}{
		"symbol":   symbol,
		"interval": interval,
		"limit":    limit,
	}

	resp, err := c.client.sendRequest("GET", endpoint, params)
	if err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}

	var response BingXResponse[[]KlineDataRaw]
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	if err := response.Error(); err != nil {
		return nil, err
	}

	var result []KlineData
	for _, data := range response.Data {
		kline, err := parseKlineData(data, interval)
		if err != nil {
			return nil, fmt.Errorf("parse: %w", err)
		}

		result = append(result, kline)
	}
	return result, nil
}

func (c *SpotClient) GetTickers(symbol ...string) (Tickers, error) {
	endpoint := "/openApi/spot/v1/ticker/24hr"
	params := map[string]interface{}{
		"timestamp": time.Now().UnixMilli(),
	}
	if len(symbol) > 0 {
		params["symbol"] = symbol[0]
	}

	resp, err := c.client.sendRequest("GET", endpoint, params)
	if err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}

	var response BingXResponse[[]TickerData]
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	if err := response.Error(); err != nil {
		return nil, err
	}

	result := Tickers{}
	for _, ticker := range response.Data {
		result[ticker.Symbol] = ticker.LastPrice
	}
	return result, nil
}
