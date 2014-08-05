package engine

type Actor interface {
	Action()
	Data() []byte
}

type board struct {
	// The top left corner of the board is (0, 0). Grows in both
	// directions.
	Width, Height int
	Actors         []Actor
}

const (
	BoardHeight = 512
	BoardWidth  = 768
)

func newBoard() *board {
	return &board{
		Width:  BoardWidth,
		Height: BoardHeight,
		Actors:  []Actor{newSheep()},
	}
}

var Board *board

func init() {
	Board = newBoard()
}

func CreateSendData() []byte {
	data := make([]byte, 0, 50)
	for _, a := range Board.Actors {
		data = append(data, a.Data()...)
	}
	return data
}

func Tick() {
	for _, a := range Board.Actors {
		a.Action()
	}
}
