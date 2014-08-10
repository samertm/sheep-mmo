package engine

import (
	"errors"

	"github.com/samertm/sheep-mmo/server/client"
	"github.com/samertm/sheep-mmo/server/message"
)

type Actor interface {
	Dataer
	Collidable
	Action()
}

type Object interface {
	Dataer
	Collidable
}

type Collidable interface {
	X() int
	Y() int
	Height() int
	Width() int
}

type Dataer interface {
	Data() []byte
}

type board struct {
	// The top left corner of the board is (0, 0). Grows in both
	// directions.
	Width, Height int
	Actors        []Actor
	Objects       []Object
}

const (
	BoardHeight = 512
	BoardWidth  = 768
)

func newBoard() *board {
	return &board{
		Width:   BoardWidth,
		Height:  BoardHeight,
		Actors:  []Actor{newSheep()},
		Objects: []Object{fence{x: 50, y: 50, width: 25, height: 25}},
	}
}

var Board *board

func (b *board) getSheep(id int) (*sheep, error) {
	for _, a := range b.Actors {
		if s, ok := a.(*sheep); ok {
			if s.id == id {
				return s, nil
			}
		}
	}
	return nil, errors.New("Could not find sheep")
}

func init() {
	Board = newBoard()
}

type mWrapper struct {
	data []byte
}

// This should never be called. @_@
func (m mWrapper) Type() string {
	return "engine"
}

func (m mWrapper) Data() []byte {
	return m.data
}

func (m mWrapper) Client() *client.C {
	return nil
}

func GenSheep() {
	Board.Actors = append(Board.Actors, newSheep())
}

func toDataerSlice(os interface{}) []Dataer {
	ds := make([]Dataer, 0)
	switch iter := os.(type) {
	case []Actor:
		for _, d := range iter {
			ds = append(ds, d)
		}
	case []Object:
		for _, d := range iter {
			ds = append(ds, d)
		}
	}
	return ds
}

func IterDataers(slices ...interface{}) <-chan Dataer {
	c := make(chan Dataer)
	go func() {
		defer close(c)
		for _, slice := range slices {
			ds := toDataerSlice(slice)
			for _, d := range ds {
				c <- d
			}
		}
	}()
	return c
}

// TODO: Rename to "Messages"?
func CreateMessages() []message.M {
	messages := make([]message.M, 0, len(Board.Actors)+len(Board.Objects))
	for d := range IterDataers(Board.Actors, Board.Objects) {
		messages = append(messages, mWrapper{data: d.Data()})
	}
	return messages
}

func Tick() {
	for _, a := range Board.Actors {
		a.Action()
	}
}
