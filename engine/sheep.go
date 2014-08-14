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
	proximateSheep     []*sheep
	proximateFlowers []*flower
	talkingTo     *sheep
	hunger        int
}

type sheepState int

const (
	thinking sheepState = iota
	startMoving
	moving
	talking
	hungry
)

func (s sheepState) String() string {
	var str string
	switch s {
	case thinking:
		str = "thinking"
	case startMoving:
		str = "startMoving"
	case moving:
		str = "moving"
	case talking:
		str = "talking"
	case hungry:
		str = "hungry"
	}
	return str
}

const (
	sheepHeight            = 40
	sheepWidth             = 38
	sheepProximateDistance = 40
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var sheepId int

func nonColliding(xrange, yrange int) (x int, y int) {
	b := box{rand.Intn(xrange), rand.Intn(yrange), sheepWidth, sheepHeight}
	// Pick a new box if there's a collision
	for collides(b, toCollidableSlice(Board.collidable), toCollidableSlice(Board.actors)) {
		b = box{rand.Intn(xrange), rand.Intn(yrange), sheepWidth, sheepHeight}
	}
	return b.x, b.y
}

func newSheep() *sheep {
	x, y := nonColliding(Board.width-sheepWidth, Board.height-sheepHeight)
	s := &sheep{
		id:        sheepId,
		x:         x,
		y:         y,
		height:    sheepHeight,
		width:     sheepWidth,
		bounceUp:  true,
		name:      "Sheepy",
		state:     thinking,
		proximateSheep: make([]*sheep, 0),
		proximateFlowers: make([]*flower, 0),
		hunger:    50,
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
	s.proximateSheep = findProximateSheep(s, Board.actors)
	s.proximateFlowers = findProximateFlowers(s, Board.noncollidable)
	switch s.state {
	case thinking:
		if rand.Intn(25) == 0 {
			s.state = startMoving
			return
		}
		if rand.Intn(15) == 0 {
			if s.hunger > 0 {
				s.hunger--
			}
		}
		if s.hunger < 45 {
			s.state = hungry
		} else if len(s.proximateSheep) != 0 && rand.Intn(10) == 0 {
			for _, sheep := range s.proximateSheep {
				if sheep.state == thinking {
					sheep.state = talking
					s.state = talking
					sheep.talkingTo = s
					s.talkingTo = sheep
					break
				}
			}
		}
	case startMoving:
		s.pickDestination()
		s.state = moving
	case moving:
		if s.arrived() {
			s.state = thinking
			return
		}
		x, y, showX, showY := s.x, s.y, s.showX, s.showY
		s.walk()
		if collides(s, toCollidableSlice(Board.collidable),
			toCollidableSlice(Board.actors)) {
			s.x, s.y, s.showX, s.showY = x, y, showX, showY
			s.state = thinking
		}
	case talking:
		if rand.Intn(50) == 0 {
			otherSheep := s.talkingTo
			otherSheep.state = thinking
			otherSheep.talkingTo = nil
			s.state = thinking
			s.talkingTo = nil
			return
		}
	case hungry:
		if len(s.proximateFlowers) != 0 {
			Board.deleteFlower(s.proximateFlowers[0].id)
			s.hunger += 10
		}
		if s.hunger >= 45 {
			s.state = thinking
			return
		}
	}
}

func findProximateSheep(s *sheep, actors []actor) []*sheep {
	distance := sheepProximateDistance // global
	proximates := make([]*sheep, 0)
	for _, a := range actors {
		if a == s {
			continue
		}
		if otherSheep, ok := a.(*sheep); ok {
			if proximate(s, otherSheep, distance) {
				proximates = append(proximates, otherSheep)
			}
		}
	}
	return proximates
}

func findProximateFlowers(s *sheep, objs []object) []*flower {
	distance := sheepProximateDistance // global
	proximates := make([]*flower, 0)
	for _, o := range objs {
		if f, ok := o.(*flower); ok {
			if proximate(s, f, distance) {
				proximates = append(proximates, f)
			}
		}
	}
	return proximates
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

// Bounding box only covers the lower half of the sheep.
func (s sheep) boundingBox() box {
	halfHeight := s.height / 2
	return box{x: s.x, y: s.y + halfHeight, width: s.width, height: halfHeight}
}

func (s sheep) data() []byte {
	return []byte(fmt.Sprintf(`(sheep %d %d %d "%s" %s)`, s.id, s.showX, s.showY, s.name, s.state))
}

func Rename(id int, name string) {
	s, err := Board.getSheep(id)
	if err != nil {
		log.Println(err)
		return
	}
	s.name = name
}
