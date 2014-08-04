package engine

import (
	"math/rand"
	"strconv"
	"time"
	"log"
)

var _ = log.Println // debugging

type Sheep struct {
	X, Y          int
	ShowX, ShowY  int
	DestX, DestY  int
	Height, Width int
	bounceHeight  int
	bounceUp      bool
	state         sheepState
}

type sheepState int

const (
	thinking sheepState = iota
	startMoving
	moving
)

type board struct {
	// The top left corner of the board is (0, 0). Grows in both
	// directions.
	Width, Height int
	Sheep         []*Sheep
}

const (
	SheepHeight = 40
	SheepWidth  = 38
	BoardHeight = 512
	BoardWidth  = 768
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
	s := &Sheep{
		X:        rand.Intn(BoardWidth - SheepWidth),
		Y:        rand.Intn(BoardHeight - SheepHeight),
		Height:   SheepHeight,
		Width:    SheepWidth,
		bounceUp: true,
		state: thinking,
	}
	s.DestX = s.X
	s.DestY = s.Y
	s.ShowX = s.X
	s.ShowY = s.Y
	return s
}

// TODO: finishMoving state, to end a bounce cleanly.
func (s *Sheep) action() {
	switch(s.state) {
	case thinking:
		//log.Println("thinking")
		if rand.Intn(25) == 0 {
			s.state = startMoving
		}
	case startMoving:
		//log.Println("start")
		s.pickDestination()
		s.state = moving
	case moving:
		//log.Println("moving")
		if s.arrived() {
			s.state = thinking
			return
		}
		s.walk()
	}
}

func (s Sheep) arrived() bool {
	return s.X == s.DestX && s.Y == s.DestY
}

func (s *Sheep) pickDestination() {
	step := 75
	s.DestX += rand.Intn(2*step) - step
	s.DestY += rand.Intn(2*step) - step
	s.correctBounds()
}

func (s *Sheep) correctBounds() {
	if s.X >= BoardWidth-SheepWidth {
		s.X = BoardWidth - SheepWidth - 1
	} else if s.X < 0 {
		s.X = 0
	}
	if s.DestX >= BoardWidth-SheepWidth {
		s.DestX = BoardWidth - SheepWidth - 1
	} else if s.DestX < 0 {
		s.DestX = 0
	}
	if s.Y >= BoardHeight-SheepHeight {
		s.Y = BoardHeight - SheepHeight - 1
	} else if s.Y < 0 {
		s.Y = 0
	}
	if s.DestY >= BoardHeight-SheepHeight {
		s.DestY = BoardHeight - SheepHeight - 1
	} else if s.DestY < 0 {
		s.DestY = 0
	}
}

func moveTowards(pos, dest, step int) int {
	if pos != dest {
		if pos < dest {
			if dest-pos < step {
				return dest
			}
			return pos + step
		}
		if pos-dest < step {
			return dest
		}
		return pos - step
	}
	return pos
}

func (s *Sheep) walk() {
	step := 5
	s.X = moveTowards(s.X, s.DestX, step)
	s.Y = moveTowards(s.Y, s.DestY, step)
	s.ShowX = s.X
	s.ShowY = s.Y
}

func (s Sheep) Data() []byte {
	return []byte("(sheep " + strconv.Itoa(s.ShowX) + " " +
		strconv.Itoa(s.ShowY) + ")")
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
		s.action()
	}
}
