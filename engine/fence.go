package engine

import "fmt"

type fence struct {
	x, y, width, height int
}

func (f fence) boundingBox() box {
	return box{x: f.x, y: f.y, width: f.width, height: f.height}
}

func (f fence) data() []byte {
	return []byte(fmt.Sprintf("(fence %d %d %d %d)", f.x, f.y, f.width, f.height))
}
