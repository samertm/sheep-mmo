package engine

import (
	"errors"
	"math"

	"github.com/samertm/sheep-mmo/server/client"
	"github.com/samertm/sheep-mmo/server/message"
)

type actor interface {
	dataer
	collidable
	action()
}

type object interface {
	dataer
	collidable
}

type collidable interface {
	boundingBox() box
}

type dataer interface {
	data() []byte
}

type board struct {
	// The top left corner of the board is (0, 0). Grows in both
	// directions.
	Width, Height int
	Actors        []actor
	Objects       []object
}

const (
	BoardHeight = 512
	BoardWidth  = 768
)

func newBoard() *board {
	return &board{
		Width:   BoardWidth,
		Height:  BoardHeight,
		Actors:  []actor{newSheep()},
		Objects: []object{fence{x: 50, y: 50, width: 25, height: 25}},
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

func toDataerSlice(os interface{}) []dataer {
	ds := make([]dataer, 0)
	switch iter := os.(type) {
	case []actor:
		for _, d := range iter {
			ds = append(ds, d)
		}
	case []object:
		for _, d := range iter {
			ds = append(ds, d)
		}
	}
	return ds
}

func IterDataers(slices ...interface{}) <-chan dataer {
	c := make(chan dataer)
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
		messages = append(messages, mWrapper{data: d.data()})
	}
	return messages
}

type pair struct {
	x, y int
}

type box struct {
	x, y, width, height int
}

func toCollidableSlice(others interface{}) []collidable {
	result := make([]collidable, 0)
	switch os := others.(type) {
	case []object:
		for _, o := range os {
			result = append(result, o)
		}
	case []actor:
		for _, o := range os {
			result = append(result, o)
		}
	}
	return result
}

func middle(b box) pair {
	return pair{int((2*b.x + b.width) / 2), int((2*b.y + b.height) / 2)}
}

func collision(c0, c1 collidable) bool {
	c0box := c0.boundingBox()
	c1box := c1.boundingBox()
	c0mid := middle(c0box)
	c1mid := middle(c1box)
	// If the distance from c0mid to c1mid is less than the distance of
	// c0mid plus half the width/height of c0 and c1, then the objects
	// intersect.
	if int(math.Abs(float64(c1mid.x - c0mid.x))) < (c0box.width + c1box.width) / 2 &&
		int(math.Abs(float64(c1mid.y - c1mid.y))) < (c0box.height + c1box.height) / 2 {
		return true
	}
	 return false
}

func collides(c collidable, cs []collidable) bool {
	for _, coll := range cs {
		if collision(c, coll) {
			return true
		}
	}
	return false
}

func Tick() {
	for _, a := range Board.Actors {
		a.action()
	}
}
