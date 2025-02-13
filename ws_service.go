package bingxgo

import (
	"encoding/json"
	"fmt"
	"strconv"

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

type Event struct {
	Code     int         `json:"code"`
	DataType string      `json:"dataType"`
	Data     interface{} `json:"data"`
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

type WsKlineEvent struct {
	Symbol    string   `json:"s"`
	Interval  Interval `json:"i"`
	Open      float64  `json:"o"`
	Close     float64  `json:"c"`
	High      float64  `json:"h"`
	Low       float64  `json:"l"`
	Volume    float64  `json:"v"`
	StartTime float64  `json:"t"`
	EndTime   float64  `json:"T"`
	Completed bool
}

type WsKlineHandler func(*WsKlineEvent)

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

	var lastEvent *WsKlineEvent

	var wsHandler = func(data []byte) {
		ev := new(Event)
		err := json.Unmarshal(data, ev)
		if err != nil {
			errHandler(err)
			return
		}

		if ev.DataType == reqEvent.DataType {
			_eventData := new(struct {
				Symbol string                   `json:"s"`
				Data   []map[string]interface{} `json:"data"`
			})
			err := json.Unmarshal(data, _eventData)
			if err != nil {
				errHandler(err)
				return
			}

			c, _ := strconv.ParseFloat(_eventData.Data[0]["c"].(string), 64)
			h, _ := strconv.ParseFloat(_eventData.Data[0]["h"].(string), 64)
			l, _ := strconv.ParseFloat(_eventData.Data[0]["l"].(string), 64)
			o, _ := strconv.ParseFloat(_eventData.Data[0]["o"].(string), 64)
			v, _ := strconv.ParseFloat(_eventData.Data[0]["v"].(string), 64)
			t := _eventData.Data[0]["T"].(float64)

			event := &WsKlineEvent{
				Symbol:    _eventData.Symbol,
				Open:      o,
				Close:     c,
				High:      h,
				Low:       l,
				Volume:    v,
				EndTime:   t,
				Completed: false,
			}

			if lastEvent == nil {
				lastEvent = event
			}

			if lastEvent.EndTime != event.EndTime {
				lastEvent.Completed = true
			}

			handler(lastEvent)

			lastEvent = event

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
	endpoint := "/openApi/user/auth/userDataStream"
	params := map[string]interface{}{}

	resp, err := c.client.sendRequest("POST", endpoint, params)
	if err != nil {
		return "", err
	}

	var response ListenKeyResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
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
