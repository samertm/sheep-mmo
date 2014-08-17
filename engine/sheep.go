package engine

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

type sheep struct {
	id               int
	x, y             int
	showX, showY     int
	destX, destY     int
	height, width    int
	bounceHeight     int
	name             string
	bounceUp         bool
	state            sheepState
	proximateSheep   []*sheep
	proximateFlowers []*flower
	talkingTo        *sheep
	hunger           int
	path             []pair
}

type sheepState int

const (
	thinking sheepState = iota
	pickDest
	moving
	talking
	hungry
	attemptToEat
)

func (s sheepState) String() string {
	var str string
	switch s {
	case thinking:
		str = "thinking"
	case pickDest:
		str = "pickDest"
	case moving:
		str = "moving"
	case talking:
		str = "talking"
	case hungry:
		str = "hungry"
	case attemptToEat:
		str = "attemptToEat"
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
		id:               sheepId,
		x:                x,
		y:                y,
		height:           sheepHeight,
		width:            sheepWidth,
		bounceUp:         true,
		name:             "Sheepy",
		state:            thinking,
		proximateSheep:   make([]*sheep, 0),
		proximateFlowers: make([]*flower, 0),
		hunger:           50,
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
			s.state = pickDest
			return
		}
		if rand.Intn(15) == 0 {
			if s.hunger > 0 {
				s.hunger--
			}
		}
		if len(s.path) != 0 {
			s.state = moving
		} else if s.hunger < 45 {
			s.state = attemptToEat
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
	case pickDest:
		s.path = []pair{s.pickDestination()}
		s.state = moving
	case moving:
		if len(s.path) == 0 {
			s.state = thinking
			return
		}
		if s.arrived(s.path[0]) {
			s.path = s.path[1:]
			return
		}
		x, y, showX, showY := s.x, s.y, s.showX, s.showY
		s.walk(s.path[0])
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
		if len(s.proximateSheep) != 0 {
			s.state = attemptToEat
			return
		}
		if s.hunger >= 45 {
			s.state = thinking
			return
		}
		for _, n := range Board.noncollidable {
			if f, ok := n.(*flower); ok {
				s.path = Board.findPath(pair{s.x, s.y}, pair{f.x, f.y})
				s.state = moving
				break

			}
		}
	case attemptToEat:
		if len(s.proximateFlowers) != 0 {
			Board.deleteFlower(s.proximateFlowers[0].id)
			s.hunger += 25
			s.state = thinking
			return
		}
		s.state = hungry
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

func (s sheep) arrived(dest pair) bool {
	blur := 10
	return s.x >= dest.x-blur && s.x <= dest.x+blur &&
		s.y >= dest.y-blur && s.y <= dest.y+blur
}

func (s sheep) pickDestination() pair {
	p := pair{s.x, s.y}
	step := 75
	p.x += rand.Intn(2*step) - step
	p.y += rand.Intn(2*step) - step
	return s.correctBounds(p)
}

func (s sheep) correctBounds(p pair) pair {
	if p.x >= Board.width-s.width {
		p.x = Board.width - s.width - 1
	} else if p.x < 0 {
		p.x = 0
	}
	if p.y >= Board.height-s.height {
		p.y = Board.height - s.height - 1
	} else if p.y < 0 {
		p.y = 0
	}
	return p
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

func (s *sheep) walk(dest pair) {
	step := 5
	s.x = moveTowards(s.x, dest.x, step)
	s.y = moveTowards(s.y, dest.y, step)
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

type pair struct {
	x, y int
}

type node struct {
	x, y int
	dist int
	up   *node
}

func minNode(queue []*node) ([]*node, *node) {
	if len(queue) == 0 {
		return nil, nil
	}
	min := math.MaxInt32
	found := false
	var foundIndex int
	for i, n := range queue {
		if n.dist >= 0 && n.dist < min {
			min = n.dist
			foundIndex = i
		}
	}
	if !found {
		return queue[1:], queue[0]
	}
	n := queue[foundIndex]
	return append(queue[:foundIndex], queue[foundIndex+1:]...), n
}

func neighbors(n *node, queue []*node) <-chan *node {
	ch := make(chan *node)
	go func() {
		defer close(ch)
		for _, inQ := range queue {
			if inQ.x >= n.x-1 && inQ.x <= n.x+1 &&
				inQ.y >= n.y-1 && inQ.y <= n.y+1 {
				ch <- inQ
			}
		}
	}()
	return ch
}

// TODO move around collidables
func (b board) findPath(start, dest pair) []pair {
	step := 30
	queue := make([]*node, 0)
	for i := 0; i < b.width/step; i++ {
		for j := 0; j < b.height/step; j++ {
			if collides(box{i * step, j * step, 10, 10},
				toCollidableSlice(Board.collidable),
				toCollidableSlice(Board.actors)) {
				continue
			}
			n := &node{x: i, y: j, dist: -1}
			if n.x == start.x/step && n.y == start.y/step {
				n.dist = 0
			}
			queue = append(queue, n)
		}
	}
	var found *node
	for len(queue) != 0 {
		var n *node
		queue, n = minNode(queue)
		if n.x == dest.x/step && n.y == dest.y/step {
			found = n
			break
		}

		for neigh := range neighbors(n, queue) {
			alt := n.dist + 1
			if neigh.dist == -1 || alt < neigh.dist {
				neigh.dist = alt
				neigh.up = n
			}
		}
	}
	var path []pair
	for on := found; on != nil; on = on.up {
		path = append([]pair{{x: on.x * step, y: on.y * step}}, path...)
	}
	return path
}
