package engine

import (
	"log"
	"math/rand"
	"strconv"
	"time"
)

type sheep struct {
	id            int
	x, y          int
	showX, showY  int
	destX, destY  int
	height, width int
	bounceHeight  int
	name          string
	bounceUp      bool
	state         sheepState
}

type sheepState int

const (
	thinking sheepState = iota
	startMoving
	moving
)

const (
	SheepHeight = 40
	SheepWidth  = 38
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var sheepId int

func newSheep() *sheep {
	s := &sheep{
		id:       sheepId,
		x:        rand.Intn(BoardWidth - SheepWidth),
		y:        rand.Intn(BoardHeight - SheepHeight),
		height:   SheepHeight,
		width:    SheepWidth,
		bounceUp: true,
		name:     "Mr. Sheep",
		state:    thinking,
	}
	sheepId++
	s.destX = s.x
	s.destY = s.y
	s.showX = s.x
	s.showY = s.y
	return s
}

// TODO: finishMoving state, to end a bounce cleanly.
func (s *sheep) Action() {
	switch s.state {
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

func (s sheep) arrived() bool {
	return s.x == s.destX && s.y == s.destY
}

func (s *sheep) pickDestination() {
	step := 75
	s.destX += rand.Intn(2*step) - step
	s.destY += rand.Intn(2*step) - step
	s.correctBounds()
}

func (s *sheep) correctBounds() {
	if s.x >= BoardWidth-SheepWidth {
		s.x = BoardWidth - SheepWidth - 1
	} else if s.x < 0 {
		s.x = 0
	}
	if s.destX >= BoardWidth-SheepWidth {
		s.destX = BoardWidth - SheepWidth - 1
	} else if s.destX < 0 {
		s.destX = 0
	}
	if s.y >= BoardHeight-SheepHeight {
		s.y = BoardHeight - SheepHeight - 1
	} else if s.y < 0 {
		s.y = 0
	}
	if s.destY >= BoardHeight-SheepHeight {
		s.destY = BoardHeight - SheepHeight - 1
	} else if s.destY < 0 {
		s.destY = 0
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

func (s *sheep) walk() {
	step := 5
	s.x = moveTowards(s.x, s.destX, step)
	s.y = moveTowards(s.y, s.destY, step)
	s.showX = s.x
	if s.bounceUp {
		s.showY += 7
	} else {
		s.showY -= 7
	}
	if s.showY > s.y+10 {
		s.bounceUp = false
	} else if s.showY < s.y {
		s.bounceUp = true
	}
}

func (s sheep) X() int {
	return s.showX
}

func (s sheep) Y() int {
	return s.showY
}

func (s sheep) Height() int {
	return s.height
}

func (s sheep) Width() int {
	return s.width
}

func (s sheep) Data() []byte {
	return []byte("(sheep " + strconv.Itoa(s.id) + " " +
		strconv.Itoa(s.showX) + " " + strconv.Itoa(s.showY) +
		" \"" + s.name + "\")")
}

func Rename(id int, name string) {
	s, err := Board.getSheep(id)
	if err != nil {
		log.Println(err)
		return
	}
	s.name = name
}
