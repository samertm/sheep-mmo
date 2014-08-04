package engine

import (
	"math/rand"
	"strconv"
	"time"
)

type Sheep struct {
	X, Y          int
	Height, Width int
}

type board struct {
	// The top left corner of the board is (0, 0). Grows in both
	// directions.
	Width, Height int
	Sheep         []*Sheep
}

const (
	SheepHeight = 40
	SheepWidth  = 38
	BoardHeight = 768
	BoardWidth  = 512
)

func newBoard() *board {
	return &board{
		Width:  BoardWidth,
		Height: BoardHeight,
		Sheep:  []*Sheep{newSheep()},
	}
}

var Board *board

func init() {
	Board = newBoard()
	rand.Seed(time.Now().UnixNano())
}

func newSheep() *Sheep {
	return &Sheep{
		X:      rand.Intn(BoardWidth - SheepWidth),
		Y:      rand.Intn(BoardHeight - SheepHeight),
		Height: SheepHeight,
		Width:  SheepWidth,
	}
}

func (s *Sheep) walk() {
	s.X = rand.Intn(BoardWidth - SheepWidth)
	s.Y = rand.Intn(BoardHeight - SheepHeight)
	
}

func (s *Sheep) Data() []byte {
	return []byte("(sheep " + strconv.Itoa(s.X) + " " + strconv.Itoa(s.Y) + ")")
}

func CreateSendData() []byte {
	data := make([]byte, 0, 50)
	for _, s := range Board.Sheep {
		data = append(data, s.Data()...)
	}
	return data
}

func Tick() {
	for _, s := range Board.Sheep {
		s.walk()
	}
}
