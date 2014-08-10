package engine

import "fmt"

type fence struct {
	x, y, width, height int
}

func (f fence) X() int {
	return f.x
}

func (f fence) Y() int {
	return f.y
}

func (f fence) Height() int {
	return f.height
}

func (f fence) Width() int {
	return f.width
}

func (f fence) Data() []byte {
	return []byte(fmt.Sprintf("(fence %d %d %d %d)", f.x, f.y, f.width, f.height))
}
