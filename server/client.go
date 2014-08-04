package server

import (
	"time"

	"github.com/gorilla/websocket"
)

func (c *client) write(messageType int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(messageType, payload)
}

func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
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
func (c *client) readPump() {
	defer func() {
		h.unregister <- c
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
		h.update <- &clientMsg{c: c, msg: message}
	}
}

func newClient(ws *websocket.Conn) *client {
	return &client{
		ws:            ws,
		send:          make(chan []byte, 256),
	}
}

type clientMsg struct {
	c   *client
	msg []byte
}

type client struct {
	ws   *websocket.Conn
	send chan []byte
}
