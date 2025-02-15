package bingxgo

import (
	"strconv"
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
		"quantity":     strconv.FormatFloat(order.Quantity, 'f', -1, 64),
		"price":        strconv.FormatFloat(order.Price, 'f', -1, 64),
	}

	resp, err := c.client.sendRequest("POST", "/openApi/swap/v2/trade/order", params)
	if err != nil {
		return nil, err
	}

	var orderResp OrderResponse
	err = json.Unmarshal(resp, &orderResp)
	return &orderResp, err
}
