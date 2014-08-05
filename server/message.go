package server

import (
	"log"
	"strconv"
)

type Message interface {
	Data() []byte
	Client() *client
}

type mouseMessage struct {
	c    *client
	x, y int
}

func (c mouseMessage) Data() []byte {
	return []byte("(mouse " + strconv.Itoa(c.x) + " " + strconv.Itoa(c.y) + ")")
}

func (c mouseMessage) Client() *client {
	return c.c
}

func decode(c *client, msg []byte) []Message {
	messages := make([]Message, 0, 2)
	ch := make(chan Message)
	go decodeRun(ch, c, msg)
	for m := range ch {
		messages = append(messages, m)
	}
	return messages
}

type stateFn func(chan Message, *client, []byte) stateFn

func decodeRun(ch chan Message, c *client, msg []byte) {
	defer close(ch)
	for state := decodeStart; state != nil; {
		state = state(ch, c, msg)
	}
}

func decodeStart(ch chan Message, c *client, msg []byte) stateFn {
	if len(msg) == 0 {
		return nil
	}
	if msg[0] == '(' {
		return decodeBeg(ch, c, msg[1:])
	}
	// Error condition.
	return nil
}

func decodeBeg(ch chan Message, c *client, msg []byte) stateFn {
	var i int
	for ; i < len(msg) && msg[i] != ' '; i++ {
		// For loop left intentionally blank.
	}
	msgType := string(msg[0:i])
	switch msgType {
	case "mouse":
		return decodeMouse(ch, c, msg[i+1:]) // i+1 skips the space
	default:
		return nil
	}
}

// msg is in the form "\d+ \d+)"
func decodeMouse(ch chan Message, c *client, msg []byte) stateFn {
	var i int
	for ; i < len(msg) && msg[i] != ' '; i++ {
		// For loop left intentionally blank.
	}
	x, err := strconv.Atoi(string(msg[0:i]))
	if err != nil {
		// Error condition.
		log.Println("errored on: " + string(msg))
		return nil
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
		return nil
	}
	ch <- mouseMessage{c: c, x: x, y: y}
	return decodeBeg
}
