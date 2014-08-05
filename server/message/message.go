package message

import (
	"log"
	"strconv"

	"github.com/samertm/sheep-mmo/server/client"
)

type M interface {
	Data() []byte
	Client() *client.C
}

type Mouse struct {
	c    *client.C
	x, y int
}

func (c Mouse) Data() []byte {
	return []byte("(mouse " + strconv.Itoa(c.x) + " " + strconv.Itoa(c.y) + ")")
}

func (c Mouse) Client() *client.C {
	return c.c
}

func Decode(c *client.C, msg []byte) []M {
	messages := make([]M, 0, 2)
	ch := make(chan M)
	go run(ch, c, msg)
	for m := range ch {
		messages = append(messages, m)
	}
	return messages
}

type stateFn func(chan M, *client.C, []byte) (stateFn, *client.C, []byte)

func run(ch chan M, c *client.C, msg []byte) {
	defer close(ch)
	for state := start; state != nil; {
		state, c, msg = state(ch, c, msg)
	}
}

func start(ch chan M, c *client.C, msg []byte) (stateFn, *client.C, []byte) {
	if len(msg) == 0 {
		return nil, nil, nil
	}
	if msg[0] == '(' {
		return beg, c, msg[1:]
	}
	// Error condition.
	return nil, nil, nil
}

func beg(ch chan M, c *client.C, msg []byte) (stateFn, *client.C, []byte) {
	var i int
	for ; i < len(msg) && msg[i] != ' '; i++ {
		// For loop left intentionally blank.
	}
	msgType := string(msg[0:i])
	switch msgType {
	case "mouse":
		return mouse, c, msg[i+1:] // i+1 skips the space
	default:
		return nil, nil, nil
	}
}

// msg is in the form "\d+ \d+)"
func mouse(ch chan M, c *client.C, msg []byte) (stateFn, *client.C, []byte) {
	var i int
	for ; i < len(msg) && msg[i] != ' '; i++ {
		// For loop left intentionally blank.
	}
	x, err := strconv.Atoi(string(msg[0:i]))
	if err != nil {
		// Error condition.
		log.Println("errored on: " + string(msg))
		return nil, nil, nil
	}
	i++ // To skip the space in msg
	var j int
	for ; j < len(msg) && msg[j] != ')'; j++ {
		// For loop left intentionally blank.
	}
	y, err := strconv.Atoi(string(msg[i:j]))
	if err != nil {
		// Error condition.
		log.Println("errored on: " + string(msg))
		return nil, nil, nil
	}
	ch <- Mouse{c: c, x: x, y: y}
	return start, c, msg[j+1:]
}
