package bingxgo

import "fmt"

// Side type of order
type SideType string

// PositionSide type of order
type PositionSideType string

// Type of order
type OrderType string

type OrderStatus string

type OrderSpecType string

type OrderWorkingType string

type Interval string

const (
	Interval1  Interval = "1m"
	Interval3  Interval = "3m"
	Interval5  Interval = "5m"
	Interval15 Interval = "15m"
	Interval30 Interval = "30m"

	Interval60  Interval = "1h"
	Interval2h  Interval = "2h"
	Interval4h  Interval = "4h"
	Interval6h  Interval = "6h"
	Interval8h  Interval = "8h"
	Interval12h Interval = "12h"

	Interval1d Interval = "1d"
	Interval3d Interval = "3d"

	Interval1w Interval = "1w"

	Interval1M Interval = "1M"
)

const (
	BuySideType  SideType = "BUY"
	SellSideType SideType = "SELL"

	ShortPositionSideType PositionSideType = "SHORT"
	LongPositionSideType  PositionSideType = "LONG"
	BothPositionSideType  PositionSideType = "BOTH"

	LimitOrderType  OrderType = "LIMIT"
	MarketOrderType OrderType = "MARKET"

	NewOrderStatus             OrderStatus = "NEW"
	PartiallyFilledOrderStatus OrderStatus = "PARTIALLY_FILLED"
	FilledOrderStatus          OrderStatus = "FILLED"
	CanceledOrderStatus        OrderStatus = "CANCELED"
	ExpiredOrderStatus         OrderStatus = "EXPIRED"

	NewOrderSpecType        OrderSpecType = "NEW"
	CanceledOrderSpecType   OrderSpecType = "CANCELED"
	CalculatedOrderSpecType OrderSpecType = "CALCULATED"
	ExpiredOrderSpecType    OrderSpecType = "EXPIRED"
	TradeOrderSpecType      OrderSpecType = "TRADE"

	MarkOrderWorkingType     OrderWorkingType = "MARK_PRICE"
	ContractOrderWorkingType OrderWorkingType = "CONTRACT_PRICE"
	IndexOrderWorkingType    OrderWorkingType = "INDEX_PRICE"
)

type BingXResponse[T any] struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	DebugMsg string `json:"debugMsg"`
	Data     T      `json:"data"`
}

func (resp BingXResponse[T]) Error() error {
	if resp.Code != 0 {
		return fmt.Errorf("code: %d, msg: %s, debugMsg: %s", resp.Code, resp.Msg, resp.DebugMsg)
	}
	return nil
}

type SpotOrderRequest struct {
	Symbol      string  `json:"symbol"`
	Side        string  `json:"side"` // BUY, SELL
	Type        string  `json:"type"` // LIMIT, MARKET
	Quantity    float64 `json:"quantity"`
	Price       float64 `json:"price,omitempty"`
	TimeInForce string  `json:"timeInForce,omitempty"` // GTC, IOC, FOK
}

type SpotOrderResponse struct {
	Symbol              string `json:"symbol"`
	OrderId             int64  `json:"orderId"`
	TransactTime        int64  `json:"transactTime"`
	Price               string `json:"price"`
	StopPrice           string `json:"stopPrice"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	Type                string `json:"type"`
	Side                string `json:"side"`
	ClientOrderID       string `json:"clientOrderID"`
}

type SpotOrder struct {
	OrderID             int     `json:"orderId"`
	ClientOrderID       string  `json:"clientOrderID"`
	Symbol              string  `json:"symbol"`
	Price               string  `json:"price"`
	OrigQty             string  `json:"origQty"`
	ExecutedQty         string  `json:"executedQty"`
	CummulativeQuoteQty string  `json:"cummulativeQuoteQty"`
	OrigQuoteQty        string  `json:"origQuoteOrderQty"`
	Status              string  `json:"status"`
	Type                string  `json:"type"`
	Side                string  `json:"side"`
	Time                int64   `json:"time"`
	UpdateTime          int64   `json:"updateTime"`
	Fee                 float64 `json:"fee"`
	FeeAsset            string  `json:"feeAsset"`
	AvgPrice            float64 `json:"avgPrice"`
}

type SpotBalance struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

type Balance struct {
	Available float64 `json:"available"`
	Locked    float64 `json:"locked"`
}

type Kline struct {
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
	Time   string  `json:"time"`
}

type OrderRequest struct {
	Symbol       string  `json:"symbol"`
	Side         string  `json:"side"`         // BUY, SELL
	PositionSide string  `json:"positionSide"` // LONG, SHORT
	Type         string  `json:"type"`         // LIMIT, MARKET
	Quantity     float64 `json:"quantity"`
	Price        float64 `json:"price"`
}

type OrderResponse struct {
	OrderId       int    `json:"orderId"`
	Symbol        string `json:"symbol"`
	Status        string `json:"status"`
	ClientOrderId string `json:"clientOrderId"`
}

type OrderBook struct {
	// The timestamp of when the orderbook last changed (in milliseconds)
	Timestamp int64 `json:"ts,omitempty"`
	// Asks order depth
	Asks [][]string `json:"asks"`
	// Bids order depth
	Bids [][]string `json:"bids"`
}

type SymbolInfos struct {
	Symbols []SymbolInfo `json:"symbols"`
}

type SymbolInfo struct {
	Symbol       string  `json:"symbol"`
	TickSize     float64 `json:"tickSize"`
	StepSize     float64 `json:"stepSize"`
	MinNotional  float64 `json:"minNotional"`
	MaxNotional  float64 `json:"maxNotional"`
	Status       int     `json:"status"`
	MaintainTime int64   `json:"maintainTime"`
}

type KlineDataRaw []float64

type KlineData struct {
	StartTime int64   `json:"startTime"`
	EndTime   int64   `json:"endTime"`
	Interval  string  `json:"interval"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
}

type TickerData struct {
	Symbol    string  `json:"symbol"`
	LastPrice float64 `json:"lastPrice"`
}

// symbol -> last price
type Tickers map[string]float64
