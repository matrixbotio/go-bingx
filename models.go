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

const (
	httpGET  = "GET"
	httpPOST = "POST"
)

const (
	endpointAccountBalance    = "/openApi/spot/v1/account/balance"
	endpointCreateOrder       = "/openApi/spot/v1/trade/order"
	endpointCreateOrdersBatch = "/openApi/spot/v1/trade/batchOrders"
	endpointGetOpenOrders     = "/openApi/spot/v1/trade/openOrders"
	endpointCancelOrder       = "/openApi/spot/v1/trade/cancel"
	endpointCancelAllOrders   = "/openApi/spot/v1/trade/cancelOpenOrders"
	endpointGetOrderData      = "/openApi/spot/v1/trade/query"
	endpointGetOrdersHistory  = "/openApi/spot/v1/trade/historyOrders"
	endpointGetOrderBook      = "/openApi/spot/v1/market/depth"
	endpointGetSymbols        = "/openApi/spot/v1/common/symbols"
	endpointGetKlinesHistory  = "/openApi/market/his/v1/kline"
	endpointGetTickers        = "/openApi/spot/v1/ticker/24hr"
	endpointGetListenKey      = "/openApi/user/auth/userDataStream"
	endpointSwapGetKlines     = "/openApi/swap/v3/quote/klines"
	endpointSwapCreateOrder   = "/openApi/swap/v2/trade/order"
)

type Interval string

const (
	Interval1   Interval = "1min"
	Interval3   Interval = "3min"
	Interval5   Interval = "5min"
	Interval15  Interval = "15min"
	Interval30  Interval = "30min"
	Interval60  Interval = "60min"
	Interval2h  Interval = "2hour"
	Interval4h  Interval = "4hour"
	Interval6h  Interval = "6hour"
	Interval8h  Interval = "8hour"
	Interval12h Interval = "12hour"
	Interval1d  Interval = "1day"
	Interval3d  Interval = "3day"
	Interval1w  Interval = "1week"
	Interval1M  Interval = "1mon"
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
	Symbol        string  `json:"symbol"`
	Side          string  `json:"side"` // BUY, SELL
	Type          string  `json:"type"` // LIMIT, MARKET
	Quantity      float64 `json:"quantity"`
	Price         float64 `json:"price,omitempty"`
	TimeInForce   string  `json:"timeInForce,omitempty"` // GTC, IOC, FOK
	ClientOrderID string  `json:"newClientOrderId"`
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
	Fee                 string  `json:"fee"`
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

	// deprecated
	MinQty float64 `json:"minQty"`
	MaxQty float64 `json:"maxQty"`
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
