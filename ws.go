package bingxgo

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

const (
	wsReadLimit = 655350
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

var wsServe = func(
	initMessage []byte,
	config *WsConfig,
	handler WsHandler,
	errHandler ErrHandler,
) (doneC, stopC chan struct{}, err error) {
	header := http.Header{}
	header.Add("Accept-Encoding", "gzip")

	c, _, err := websocket.DefaultDialer.Dial(config.Endpoint, header)
	if err != nil {
		return nil, nil, err
	}

	if initMessage != nil {
		err = c.WriteMessage(websocket.TextMessage, initMessage)
		if err != nil {
			return nil, nil, err
		}
	}

	c.SetReadLimit(wsReadLimit)
	doneC = make(chan struct{})
	stopC = make(chan struct{})
	go func() {
		defer close(doneC)
		silent := false
		go func() {
			select {
			case <-stopC:
				silent = true
			case <-doneC:
			}
			c.Close()
		}()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if !silent {
					errHandler(err)
				}
				return
			}

			decodedMsg, err := DecodeGzip(message)
			if err != nil {
				if !silent {
					errHandler(err)
				}
				return
			}

			isPing, err := handlePing(c, decodedMsg)
			if err != nil {
				if !silent {
					errHandler(err)
				}
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
