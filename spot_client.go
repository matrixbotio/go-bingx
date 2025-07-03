package bingxgo

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

type SpotClient struct {
	client *Client
}

func NewSpotClient(client *Client) SpotClient {
	return SpotClient{client: client}
}

func (c *SpotClient) get(
	method string,
	params map[string]any,
	resultPointer any,
) error {
	return c.client.sendRequest(http.MethodGet, method, params, resultPointer)
}

func (c *SpotClient) post(
	method string,
	params map[string]any,
	resultPointer any,
) error {
	return c.client.sendRequest(http.MethodPost, method, params, resultPointer)
}

func (c *SpotClient) GetBalance() ([]SpotBalance, error) {
	params := map[string]any{
		"timestamp": time.Now().UnixMilli(),
	}

	var response BingXResponse[map[string][]SpotBalance]
	if err := c.get(endpointAccountBalance, params, &response); err != nil {
		return nil, err
	}
	if err := response.Error(); err != nil {
		return nil, err
	}

	if len(response.Data) == 0 {
		return nil, nil
	}
	return response.Data["balances"], nil
}

func (c *SpotClient) CreateOrder(order SpotOrderRequest) (*SpotOrderResponse, error) {
	params := map[string]any{
		"symbol":   order.Symbol,
		"side":     string(order.Side),
		"type":     string(order.Type),
		"quantity": decimal.NewFromFloat(order.Quantity).String(),
		"price":    decimal.NewFromFloat(order.Price).String(),
	}
	if order.ClientOrderID != "" {
		params["newClientOrderId"] = order.ClientOrderID
	}

	var response BingXResponse[SpotOrderResponse]
	if err := c.post(endpointCreateOrder, params, &response); err != nil {
		return nil, err
	}

	if err := response.Error(); err != nil {
		return nil, err
	}
	return &response.Data, nil
}

func (c *SpotClient) CreateBatchOrders(
	orders []SpotOrderRequest,
	isSync bool,
) ([]SpotOrderResponse, error) {
	ordersJSON, err := json.Marshal(orders)
	if err != nil {
		return nil, err
	}
	params := map[string]any{
		"data": string(ordersJSON),
		"sync": isSync,
	}

	var response BingXResponse[map[string][]SpotOrderResponse]
	if err := c.post(endpointCreateOrdersBatch, params, &response); err != nil {
		return nil, err
	}

	if err := response.Error(); err != nil {
		return nil, err
	}
	return response.Data["orders"], err
}

func (c *SpotClient) GetOpenOrders(symbol string) ([]SpotOrder, error) {
	params := map[string]any{
		"symbol": symbol,
	}

	var response BingXResponse[map[string][]SpotOrder]
	if err := c.get(endpointGetOpenOrders, params, &response); err != nil {
		return nil, err
	}

	if err := response.Error(); err != nil {
		return nil, err
	}
	return response.Data["orders"], nil
}

func (c *SpotClient) CancelOrder(symbol string, orderId string) error {
	params := map[string]any{
		"symbol":  symbol,
		"orderId": orderId,
	}

	var response BingXResponse[any]
	if err := c.post(endpointCancelOrder, params, &response); err != nil {
		return err
	}

	return response.Error()
}

func (c *SpotClient) CancelOrderByClientOrderID(
	symbol string,
	clientOrderID string,
) error {
	params := map[string]any{
		"symbol":        symbol,
		"clientOrderID": clientOrderID,
	}

	var response BingXResponse[any]
	if err := c.post(endpointCancelOrder, params, &response); err != nil {
		return err
	}

	return response.Error()
}

func (c *SpotClient) CancelAllOpenOrders(symbol string) error {
	params := map[string]any{
		"symbol": symbol,
	}

	var response BingXResponse[any]
	if err := c.post(endpointCancelAllOrders, params, &response); err != nil {
		return err
	}

	return response.Error()
}

func (c *SpotClient) GetOrder(symbol string, orderID int64) (*SpotOrder, error) {
	return c.getOrderData(map[string]any{
		"symbol":    symbol,
		"orderId":   orderID,
		"timestamp": time.Now().UnixMilli(),
	})
}

func (c *SpotClient) GetOrderByClientOrderID(
	symbol string,
	clientOrderID string,
) (*SpotOrder, error) {
	return c.getOrderData(map[string]any{
		"symbol":        symbol,
		"clientOrderID": clientOrderID,
		"timestamp":     time.Now().UnixMilli(),
	})
}

func (c *SpotClient) getOrderData(
	params map[string]any,
) (*SpotOrder, error) {
	var response BingXResponse[SpotOrder]
	if err := c.get(endpointGetOrderData, params, &response); err != nil {
		return nil, err
	}

	if err := response.Error(); err != nil {
		return nil, err
	}
	return &response.Data, nil
}

type OrdersHistoryResponse struct {
	Orders []HistoryOrder `json:"orders"`
}

// GetHistoryOrder - get order from trade history. Time in unix timestamp ms
func (c *SpotClient) GetHistoryOrder(
	symbol string,
	orderID int64,
) (HistoryOrder, error) {
	orders, err := c.HistoryOrders(symbol, 0, 0, orderID)
	if err != nil {
		return HistoryOrder{}, fmt.Errorf("get history: %w", err)
	}

	for _, order := range orders {
		if order.OrderID == orderID {
			return order, nil
		}
	}
	return HistoryOrder{}, errors.New("not found")
}

// HistoryOrders - get orders history. Time in unix timestamp ms
func (c *SpotClient) HistoryOrders(
	symbol string,
	fromTime int64,
	toTime int64,
	orderID int64,
) ([]HistoryOrder, error) {
	params := map[string]any{
		"symbol":    symbol,
		"timestamp": time.Now().UnixMilli(),
	}

	if orderID > 0 {
		params["orderId"] = orderID
	}
	if fromTime > 0 {
		params["startTime"] = fromTime
	}
	if toTime > 0 {
		params["endTime"] = toTime
	}

	var response BingXResponse[OrdersHistoryResponse]
	if err := c.get(endpointGetOrdersHistory, params, &response); err != nil {
		return nil, err
	}

	if err := response.Error(); err != nil {
		return nil, err
	}

	return response.Data.Orders, nil
}

func (c *SpotClient) OrderBook(symbol string, limit int) (*OrderBook, error) {
	params := map[string]any{
		"symbol":    symbol,
		"timestamp": time.Now().UnixMilli(),
	}
	if limit > 0 {
		params["limit"] = limit
	}

	var response BingXResponse[OrderBook]
	if err := c.get(endpointGetOrderBook, params, &response); err != nil {
		return nil, err
	}

	if err := response.Error(); err != nil {
		return nil, err
	}
	return &response.Data, nil
}

func (c *SpotClient) GetSymbols(symbol ...string) ([]SymbolInfo, error) {
	params := map[string]any{
		"timestamp": time.Now().UnixMilli(),
	}
	if len(symbol) > 0 {
		params["symbol"] = symbol[0]
	}

	var response BingXResponse[SymbolInfos]
	if err := c.get(endpointGetSymbols, params, &response); err != nil {
		return nil, err
	}

	if err := response.Error(); err != nil {
		return nil, err
	}
	return response.Data.Symbols, nil
}

func (c *SpotClient) GetHistoricalCandles(
	symbol string,
	interval string,
	limit int64,
) ([]KlineData, error) {
	params := map[string]any{
		"symbol":   symbol,
		"interval": interval,
		"limit":    limit,
	}

	var response BingXResponse[[]KlineDataRaw]
	if err := c.get(endpointGetCandlesHistory, params, &response); err != nil {
		return nil, fmt.Errorf("send: %w", err)
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

func (c *SpotClient) GetHistoricalKlines(
	symbol string,
	interval string,
	limit int64,
) ([]KlineData, error) {
	params := map[string]any{
		"symbol":   symbol,
		"interval": interval,
		"limit":    limit,
		"endTime":  time.Now().UnixMilli(),
	}

	var response BingXResponse[[]KlineDataRaw]
	if err := c.get(endpointGetKlinesHistory, params, &response); err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}

	if err := response.Error(); err != nil {
		return nil, err
	}
	if response.Msg != "" {
		return nil, errors.New(response.Msg)
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
	params := map[string]any{
		"timestamp": time.Now().UnixMilli(),
	}
	if len(symbol) > 0 {
		params["symbol"] = symbol[0]
	}

	var response BingXResponse[[]TickerData]
	if err := c.get(endpointGetTickers, params, &response); err != nil {
		return nil, fmt.Errorf("send: %w", err)
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
