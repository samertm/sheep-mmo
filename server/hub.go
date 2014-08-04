package server

import (
	"log"
	"time"

	"github.com/samertm/sheep-mmo/engine"
)

type hub struct {
	clients    map[*client]int
	broadcast  chan []byte
	register   chan *client
	unregister chan *client
	update     chan *clientMsg
	tick       <-chan time.Time
}

var h = hub{
	clients:    make(map[*client]int),
	broadcast:  make(chan []byte),
	register:   make(chan *client),
	unregister: make(chan *client),
	update:     make(chan *clientMsg),
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
			close(c.send)
		case cMsg := <-h.update: // recieves clientMsg
			log.Println(string(cMsg.msg))
		case <-h.tick:
			if len(h.clients) != 0 {
				go func() {
					engine.Tick()
					h.broadcast <- engine.CreateSendData()
				}()
			}
		case d := <-h.broadcast:
			for c := range h.clients {
				c.send <- d
			}
		}
	}
}
