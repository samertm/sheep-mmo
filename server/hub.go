package server

import (
	"time"

	"github.com/samertm/sheep-mmo/engine"
	"github.com/samertm/sheep-mmo/server/client"
	"github.com/samertm/sheep-mmo/server/message"
)

type hub struct {
	clients    map[*client.C]*data
	broadcast  chan message.M
	register   chan *client.C
	unregister chan *client.C
	update     chan *client.CMsg
	tick       <-chan time.Time
}

var h = hub{
	clients:    make(map[*client.C]*data),
	broadcast:  make(chan message.M),
	register:   make(chan *client.C),
	unregister: make(chan *client.C),
	update:     make(chan *client.CMsg),
	tick:       time.Tick(70 * time.Millisecond),
}

var forclient map[string]bool

func init() {
	forclient = map[string]bool{
		"mouse": true,
	}
}

func (h *hub) run() {
	id := 0
	for {
		select {
		case c := <-h.register:
			h.clients[c] = &data{
				id:   id,
				msgs: make(map[string]message.M),
			}
			id++
		case c := <-h.unregister:
			delete(h.clients, c)
			close(c.Send)
		case cMsg := <-h.update: // recieves clientMsg
			msgs := message.Decode(cMsg.C, cMsg.Msg)
			for _, m := range msgs {
				if _, ok := forclient[m.Type()]; ok {
					h.clients[cMsg.C].msgs[m.Type()] = m
					go func(m message.M) {
						h.broadcast <- m
					}(m)
				} else if m.Type() == "rename" {
					r := m.(message.Rename)
					engine.Rename(r.Id, r.Name)
				} else if m.Type() == "gen-sheep" {
					engine.GenSheep()
				}
			}
		case <-h.tick:
			if len(h.clients) != 0 {
				engine.Tick()
				msgs := message.TickConcat(engine.CreateMessages(),
					createMessages(h.clients))
				for _, m := range msgs {
					go func(m message.M) {
						h.broadcast <- m
					}(m)
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
