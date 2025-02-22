package bingxgo

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

const (
	wsReadLimit            = 655350
	extendListenKeyTimeout = time.Minute * 40
	pingListenKeyTimeout   = time.Minute * 30
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// WsHandler handle raw websocket message
type WsHandler func([]byte)

// ErrHandler handles errors
type ErrHandler func(err error)

// WsConfig webservice configuration
type WsConfig struct {
	Endpoint string
}

type PingMessage struct {
	ID   string `json:"ping"`
	Time string `json:"time"`
}

type PongMessage struct {
	ID   string `json:"pong"`
	Time string `json:"time"`
}

func newWsConfig(endpoint string) *WsConfig {
	return &WsConfig{
		Endpoint: endpoint,
	}
}

func extendListenKey(client *Client, keyID string) error {
	var result any
	return client.sendRequest(
		http.MethodPut,
		endpointExtendListenKey,
		map[string]interface{}{"listenKey": keyID},
		&result,
	)
}

func wsServe(
	initMessage []byte,
	client *Client,
	listenKeyID string,
	config *WsConfig,
	handler WsHandler,
	errHandler ErrHandler,
) (doneC, stopC chan struct{}, err error) {
	header := http.Header{}
	header.Add("Accept-Encoding", "gzip")

	wsConn, _, err := websocket.DefaultDialer.Dial(config.Endpoint, header)
	if err != nil {
		return nil, nil, err
	}

	if initMessage != nil {
		err = wsConn.WriteMessage(websocket.TextMessage, initMessage)
		if err != nil {
			return nil, nil, err
		}
	}

	wsConn.SetReadLimit(wsReadLimit)
	doneC = make(chan struct{})
	stopC = make(chan struct{})
	active := true

	go func() {
		defer close(doneC)
		silent := false

		// await stop
		go func() {
			select {
			case <-stopC:
				silent = true
			case <-doneC:
			}
			active = false
			wsConn.Close()
		}()

		if listenKeyID != "" {
			// auto-extend listen-key
			go func() {
				for {
					time.Sleep(extendListenKeyTimeout)

					if !active {
						return
					}

					if err := extendListenKey(client, listenKeyID); err != nil {
						if !silent {
							errHandler(fmt.Errorf("extend listen key: %w", err))
						}
					}
				}
			}()

			// periodicaly ping listen-key
			go func() {
				for {
					time.Sleep(pingListenKeyTimeout)

					if !active {
						return
					}

					if err := wsConn.WriteMessage(
						websocket.PingMessage,
						[]byte("keepalive"),
					); err != nil && !silent {
						errHandler(err)
					}
				}
			}()
		}

		// read messages
		for {
			_, message, err := wsConn.ReadMessage()
			if err != nil {
				if !silent {
					errHandler(err)
				}
				return
			}

			decodedMsg, err := DecodeGzip(message)
			if err != nil && !silent {
				errHandler(err)
				return
			}

			isPing, err := handlePing(wsConn, decodedMsg)
			if err != nil && !silent {
				errHandler(err)
				continue
			}
			if isPing {
				continue
			}

			handler(decodedMsg)
		}
	}()
	return
}

func handlePing(c *websocket.Conn, msg []byte) (bool, error) {
	if msg == nil {
		return false, nil
	}

	if !strings.Contains(string(msg), `"ping"`) {
		return false, nil
	}

	var pingMsg PingMessage
	if err := json.Unmarshal(msg, &pingMsg); err != nil {
		return false, fmt.Errorf("decode ping message: %w", err)
	}

	if pingMsg.ID == "" {
		return false, nil
	}

	if err := pong(c, pingMsg); err != nil {
		return false, fmt.Errorf("pong: %w", err)
	}
	return true, nil
}

func pong(c *websocket.Conn, pingMsg PingMessage) error {
	msg := PongMessage{
		ID:   pingMsg.ID,
		Time: pingMsg.Time,
	}

	if err := c.WriteJSON(msg); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
