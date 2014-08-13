package engine

import "fmt"

type flower struct {
	id                  int
	x, y, width, height int
}

var flowerId int

func newFlower(x, y int) *flower {
	f := &flower{
		x:      x,
		y:      y,
		width:  15,
		height: 15,
	}
	if Board.collisions(f) {
		return nil
	}
	f.id = flowerId
	flowerId++
	return f
}

func (f flower) boundingBox() box {
	return box{x: f.x, y: f.y, width: f.width, height: f.height}
}

func (f flower) data() []byte {
	return []byte(fmt.Sprintf("(flower %d %d %d)", f.id, f.x, f.y))
}
