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
	width, height int
	actors        []actor
	objects       []object
}

const (
	boardHeight = 768
	boardWidth = 512
)

func newBoard() *board {
	return &board{
		width:   boardWidth,
		height:  boardHeight,
		actors:  make([]actor, 0, 1),
		objects: make([]object, 0, 1),
	}
}

var Board *board

func (b *board) getSheep(id int) (*sheep, error) {
	for _, a := range b.actors {
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
	Board.actors = append(Board.actors, newSheep())
	Board.objects = append(Board.objects, fence{50,50,25,25})
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
	Board.actors = append(Board.actors, newSheep())
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
	messages := make([]message.M, 0, len(Board.actors)+len(Board.objects))
	for d := range IterDataers(Board.actors, Board.objects) {
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

func proximate(c0, c1 collidable, distance int) bool {
	c0box := c0.boundingBox()
	c1box := c1.boundingBox()
	c0mid := middle(c0box)
	c1mid := middle(c1box)
	// If the distance from c0mid to c1mid is less than the distance of
	// c0mid plus half the width/height of c0 and c1, then the objects
	// intersect. Add `distance' to the second number, so that we can
	// adjust how close we can be to this other object before we 
	// determinte that we are proximate to it.
	if int(math.Abs(float64(c1mid.x-c0mid.x))) < (c0box.width+c1box.width)/2 + distance &&
		int(math.Abs(float64(c1mid.y-c1mid.y))) < (c0box.height+c1box.height)/2 + distance {
		return true
	}
	return false
}

func collision(c0, c1 collidable) bool {
	return proximate(c0, c1, 0)
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
	for _, a := range Board.actors {
		a.action()
	}
}
