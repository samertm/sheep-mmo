package client

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func (c *C) write(messageType int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(messageType, payload)
}

func (c *C) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// readPump pumps messages from the websocket connection to the hub.
func (c *C) ReadPump(unregister chan *C, update chan *CMsg) {
	defer func() {
		unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		update <- &CMsg{C: c, Msg: message}
	}
}

func New(ws *websocket.Conn) *C {
	return &C{
		ws:   ws,
		Send: make(chan []byte, 256),
	}
}

type CMsg struct {
	C   *C
	Msg []byte
}

type C struct {
	ws   *websocket.Conn
	Send chan []byte
}
