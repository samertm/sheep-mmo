package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/samertm/sheep-mmo/server/client"
	"github.com/samertm/sheep-mmo/server/message"
)

type data struct {
	id   int
	msgs map[string]message.M
}

func createMessages(cs map[*client.C]*data) []message.M {
	msgs := make([]message.M, 0)
	for _, d := range cs {
		for _, m := range d.msgs {
			msgs = append(msgs, m)
		}
	}
	return msgs
}

func handleWs(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	c := client.New(ws)
	h.register <- c
	go c.WritePump()
	c.ReadPump(h.unregister, h.update)
}

func ListenAndServe(ipaddr string) {
	http.HandleFunc("/", handleWs)
	go h.run()
	if err := http.ListenAndServe(ipaddr, nil); err != nil {
		log.Fatal(err)
	}
}
