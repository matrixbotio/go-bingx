package bingxgo

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	baseWsUrl        = "wss://open-api-ws.bingx.com/market"
	baseAccountWsUrl = "wss://open-api-ws.bingx.com/market?listenKey="
)

func getWsEndpoint() string {
	return baseWsUrl
}

func getAccountWsEndpoint(listenKey string) string {
	return baseAccountWsUrl + listenKey
}

type Event[dataType any] struct {
	Code     int      `json:"code"`
	DataType string   `json:"dataType"`
	Data     dataType `json:"data"`
}

type KlineEventData struct {
	EventTime  int64      `json:"E"`
	EventType  string     `json:"e"`
	PairSymbol string     `json:"s"`
	Kline      KlineEvent `json:"K"`
}

type WsRequestType string

const (
	SubscribeRequestType  WsRequestType = "sub"
	UnubscribeRequestType WsRequestType = "unsub"
)

type RequestEvent struct {
	Id       uuid.UUID     `json:"id"`
	ReqType  WsRequestType `json:"reqType"`
	DataType string        `json:"dataType"`
}

type KlineEvent struct {
	Symbol    string   `json:"s"`
	Interval  Interval `json:"i"`
	Open      string   `json:"o"`
	Close     string   `json:"c"`
	High      string   `json:"h"`
	Low       string   `json:"l"`
	Volume    string   `json:"v"`
	StartTime int64    `json:"t"`
	EndTime   int64    `json:"T"`

	Completed bool  `json:"completed"`
	EventTime int64 `json:"eventTime"`
}

type WsKlineHandler func(KlineEvent)

func WsKlineServe(
	symbol string,
	interval Interval,
	handler WsKlineHandler,
	errHandler ErrHandler,
) (doneC, stopC chan struct{}, err error) {
	// Symbol e.g. "BTC-USDT"
	// Interval e.g. "1m", "3h"
	reqEvent := RequestEvent{
		Id:       uuid.New(),
		ReqType:  SubscribeRequestType,
		DataType: fmt.Sprintf("%s@kline_%s", symbol, interval),
	}

	var lastEventEndTime int64

	var wsHandler = func(data []byte) {
		if data == nil {
			return
		}

		if strings.Contains(string(data), "error: ") {
			errHandler(errors.New(string(data)))
			return
		}

		var ev Event[KlineEventData]
		err := json.Unmarshal(data, &ev)
		if err != nil {
			errHandler(err)
			return
		}

		if ev.DataType == reqEvent.DataType {
			if lastEventEndTime == 0 {
				lastEventEndTime = ev.Data.Kline.EndTime
			}

			if lastEventEndTime != ev.Data.Kline.EndTime {
				lastEventEndTime = ev.Data.Kline.EndTime
				ev.Data.Kline.Completed = true
			}

			ev.Data.Kline.EventTime = ev.Data.EventTime
			handler(ev.Data.Kline)
		}
	}

	initMessage, err := json.Marshal(reqEvent)
	if err != nil {
		return nil, nil, err
	}

	return wsServe(initMessage, newWsConfig(getWsEndpoint()), wsHandler, errHandler)
}

type WsOrder struct {
	TransactionID string        `json:"t"`
	Symbol        string        `json:"s"`
	Side          SideType      `json:"S"`
	OrderType     OrderType     `json:"o"`
	Price         string        `json:"p"`
	AveragePrice  string        `json:"ap"`
	Quantity      string        `json:"q"`
	Amount        string        `json:"Q"`
	StopPrice     string        `json:"sp"`
	Status        OrderStatus   `json:"X"`
	EventType     OrderSpecType `json:"x"`
	Timestamp     int           `json:"T"`
	OrderID       int           `json:"i"`
	ClientOrderID string        `json:"c"`
}

type WsOrderUpdateEvent struct {
	EventType string   `json:"e"`
	Time      int      `json:"E"`
	Order     *WsOrder `json:"o"`
}

type WsOrderUpdateHandler func(*WsOrder)

type ListenKeyResponse struct {
	Key string `json:"listenKey"`
}

func (c *SpotClient) getListenKey() (string, error) {
	var response ListenKeyResponse
	if err := c.post(
		endpointGetListenKey,
		map[string]interface{}{},
		&response,
	); err != nil {
		return "", err
	}

	return response.Key, nil
}

func (c *SpotClient) WsOrderUpdateServe(
	handler WsOrderUpdateHandler,
	errHandler ErrHandler,
) (doneC, stopC chan struct{}, err error) {
	listenKey, err := c.getListenKey()
	if err != nil {
		return nil, nil, fmt.Errorf("get listen key: %w", err)
	}

	var wsHandler = func(data []byte) {

		var evMap map[string]interface{}
		err := json.Unmarshal(data, &evMap)
		if err != nil {
			errHandler(err)
			return
		}

		if evMap["e"].(string) == "ORDER_TRADE_UPDATE" {
			event := new(WsOrderUpdateEvent)
			err = json.Unmarshal(data, event)
			if err != nil {
				errHandler(err)
				return
			}
			handler(event.Order)
		}
	}

	return wsServe(
		nil,
		newWsConfig(getAccountWsEndpoint(listenKey)),
		wsHandler, errHandler,
	)
}
