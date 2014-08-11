package engine

import (
	"fmt"
	"log"
	"math/rand"
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
	sheepHeight = 40
	sheepWidth  = 38
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var sheepId int

func newSheep() *sheep {
	s := &sheep{
		id:       sheepId,
		x:        rand.Intn(Board.width - sheepWidth),
		y:        rand.Intn(Board.height - sheepHeight),
		height:   sheepHeight,
		width:    sheepWidth,
		bounceUp: true,
		name:     "Sheepy",
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
func (s *sheep) action() {
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
		x, y, showX, showY := s.x, s.y, s.showX, s.showY
		s.walk()
		if collides(s, toCollidableSlice(Board.objects)) {
			s.x, s.y, s.showX, s.showY = x, y, showX, showY
			s.state = thinking
		}
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
	if s.x >= Board.width-s.width {
		s.x = Board.width - s.width - 1
	} else if s.x < 0 {
		s.x = 0
	}
	if s.destX >= Board.width-s.width {
		s.destX = Board.width - s.width - 1
	} else if s.destX < 0 {
		s.destX = 0
	}
	if s.y >= Board.height-s.height {
		s.y = Board.height - s.height - 1
	} else if s.y < 0 {
		s.y = 0
	}
	if s.destY >= Board.height-s.height {
		s.destY = Board.height - s.height - 1
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

func (s sheep) boundingBox() box {
	return box{x: s.x, y: s.y, width: s.width, height: s.height}
}

func (s sheep) data() []byte {
	return []byte(fmt.Sprintf(`(sheep %d %d %d "%s")`, s.id, s.showX, s.showY, s.name))
}

func Rename(id int, name string) {
	s, err := Board.getSheep(id)
	if err != nil {
		log.Println(err)
		return
	}
	s.name = name
}
