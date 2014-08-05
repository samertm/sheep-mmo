package engine

import (
	"github.com/samertm/sheep-mmo/server/client"
	"github.com/samertm/sheep-mmo/server/message"
)

type Actor interface {
	Action()
	Data() []byte
}

type board struct {
	// The top left corner of the board is (0, 0). Grows in both
	// directions.
	Width, Height int
	Actors        []Actor
}

const (
	BoardHeight = 512
	BoardWidth  = 768
)

func newBoard() *board {
	return &board{
		Width:  BoardWidth,
		Height: BoardHeight,
		Actors: []Actor{newSheep()},
	}
}

var Board *board

func init() {
	Board = newBoard()
}

type mWrapper struct {
	data []byte
}

func (m mWrapper) Data() []byte {
	return m.data
}

func (m mWrapper) Client() *client.C {
	return nil
}

// TODO: Rename to "Messages"?
func CreateMessages() []message.M {
	messages := make([]message.M, 0, len(Board.Actors))
	for _, a := range Board.Actors {
		messages = append(messages, mWrapper{data: a.Data()})
	}
	return messages
}

func Tick() {
	for _, a := range Board.Actors {
		a.Action()
	}
}
