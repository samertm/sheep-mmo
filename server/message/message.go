package message

import (
	"fmt"
	"log"
	"strconv"

	"github.com/samertm/sheep-mmo/server/client"
)

type M interface {
	Type() string
	Data() []byte
	Client() *client.C
}

var mouseIds map[*client.C]int
var maxMouseId int

func init() {
	mouseIds = make(map[*client.C]int)
}

type tick struct {
}

func (t tick) Data() []byte {
	return []byte("(tick)")
}

func (t tick) Client() *client.C {
	return nil
}

func (t tick) Type() string {
	return "tick"
}

func TickConcat(msgs ...[]M) []M {
	result := []M{tick{}}
	for _, m := range msgs {
		result = append(result, m...)
	}
	return result
}

type Mouse struct {
	Id   int
	c    *client.C
	X, Y int
}

func NewMouse(c *client.C, x, y int) Mouse {
	var id int
	if i, ok := mouseIds[c]; ok {
		id = i
	} else {
		id = maxMouseId
		mouseIds[c] = id
		maxMouseId++
	}
	return Mouse{Id: id, c: c, X: x, Y: y}
}

func (m Mouse) Data() []byte {
	return []byte(fmt.Sprintf("(mouse %d %d %d)", m.Id, m.X, m.Y))
}

func (m Mouse) Client() *client.C {
	return m.c
}

func (m Mouse) Type() string {
	return "mouse"
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
	for ; i < len(msg) && msg[i] != ' ' && msg[i] != ')'; i++ {
		// For loop left intentionally blank.
	}
	msgType := string(msg[0:i])
	switch msgType {
	case "mouse":
		return mouseMsg, c, msg[i+1:] // i+1 skips the space
	case "rename":
		return renameMsg, c, msg[i+1:]
	case "gen-sheep":
		return gensheepMsg, c, msg[i+1:]
	default:
		log.Println("Errored: " + string(msg))
		return nil, nil, nil
	}
}

// msg is in the form "\d+ \d+)"
func mouseMsg(ch chan M, c *client.C, msg []byte) (stateFn, *client.C, []byte) {
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
	ch <- NewMouse(c, x, y)
	return start, c, msg[j+1:]
}

// Returns (stringfound, restofmessage)
func getString(msg []byte) (string, []byte) {
	if len(msg) == 0 {
		return "", msg
	}
	for i := 0; i < len(msg); i++ {
		if msg[i] == '"' {
			return string(msg[:i]), msg[i+1:]
		}
	}
	return "", msg
}

type Rename struct {
	Id   int
	Name string
}

func (r Rename) Data() []byte {
	return []byte("")
}

func (r Rename) Client() *client.C {
	return nil
}

func (r Rename) Type() string {
	return "rename"
}

// (rename \d \w)
func renameMsg(ch chan M, c *client.C, msg []byte) (stateFn, *client.C, []byte) {
	var i int
	for ; i < len(msg) && msg[i] != ' '; i++ {
	}
	id, err := strconv.Atoi(string(msg[:i]))
	if err != nil {
		log.Println("errored on: " + string(msg))
		return nil, nil, nil
	}
	i++ // skip space
	var name string
	if msg[i] == '"' {
		name, msg = getString(msg[i+1:])
	} else {
		j := i
		for ; j < len(msg) && msg[j] != ')'; j++ {
		}
		name, msg = string(msg[i:j]), msg[j:]
	}
	if msg[0] != ')' {
		// TODO tighten up
		log.Println("errored " + string(msg))
		return nil, nil, nil
	}
	ch <- Rename{Id: id, Name: name}
	return start, c, msg[1:]
}

type GenSheep struct {
}

func (g GenSheep) Data() []byte {
	return []byte("")
}

func (g GenSheep) Client() *client.C {
	return nil
}

func (g GenSheep) Type() string {
	return "gen-sheep"
}

func gensheepMsg(ch chan M, c *client.C, msg []byte) (stateFn, *client.C, []byte) {
	ch <- GenSheep{}
	var i int
	for ; i < len(msg) && msg[i] != ')'; i++ {
	}
	// TODO: fix this ):L
	return start, c, msg
}
