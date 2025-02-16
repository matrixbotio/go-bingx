package bingxgo

import (
	"github.com/shopspring/decimal"
)

type TradeClient struct {
	client *Client
}

func (c *TradeClient) CreateOrder(order OrderRequest) (*OrderResponse, error) {
	params := map[string]interface{}{
		"symbol":       order.Symbol,
		"side":         string(order.Side),
		"positionSide": string(order.PositionSide),
		"type":         string(order.Type),
		"quantity":     decimal.NewFromFloat(order.Quantity).String(),
		"price":        decimal.NewFromFloat(order.Price).String(),
	}

	resp, err := c.client.post(endpointSwapCreateOrder, params)
	if err != nil {
		return nil, err
	}

	var orderResp OrderResponse
	err = json.Unmarshal(resp, &orderResp)
	return &orderResp, err
}
