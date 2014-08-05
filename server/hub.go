package server

import (
	"log"
	"time"

	"github.com/samertm/sheep-mmo/engine"
	"github.com/samertm/sheep-mmo/server/client"
	"github.com/samertm/sheep-mmo/server/message"
)

type hub struct {
	clients    map[*client.C]int
	broadcast  chan message.M
	register   chan *client.C
	unregister chan *client.C
	update     chan *client.CMsg
	tick       <-chan time.Time
}

var h = hub{
	clients:    make(map[*client.C]int),
	broadcast:  make(chan message.M),
	register:   make(chan *client.C),
	unregister: make(chan *client.C),
	update:     make(chan *client.CMsg),
	tick:       time.Tick(70 * time.Millisecond),
}

func (h *hub) run() {
	id := 0
	for {
		select {
		case c := <-h.register:
			h.clients[c] = id
			id++
		case c := <-h.unregister:
			delete(h.clients, c)
			close(c.Send)
		case cMsg := <-h.update: // recieves clientMsg
			msgs := message.Decode(cMsg.C, cMsg.Msg)
			for _, m := range msgs {
				go func() {
					h.broadcast <- m
				}()
			}
		case <-h.tick:
			if len(h.clients) != 0 {
				engine.Tick()
				msgs := engine.CreateMessages()
				for _, m := range msgs {
					go func() {
						h.broadcast <- m
					}()
				}
			}
		case m := <-h.broadcast:
			for c := range h.clients {
				if c != m.Client() {
					c.Send <- m.Data()
				}
			}
		}
	}
}
