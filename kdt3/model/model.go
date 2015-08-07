package model

type Player struct {
        PlayerId string
        GameId string
        Handle string
}

type Game struct {
        GameId string
        PlayerIds []string
        Owner int
        Turn int
        Wins []Point
        Board *Board
}

type Board struct {
        K int
        D *Cell
        Size int
}

func NewBoard(K, size int) *Board {
        board := &Board{K:K, Size: size, D: &Cell{}}
        var populate func(*Cell, int)
        populate = func(cell *Cell, depth int) {
                if depth > 0 {
                        cell.D = make([]*Cell, size)
                        for i, _ := range cell.D {
                                cell.D[i] = &Cell{}
                                populate(cell.D[i], depth -1)
                        }
                }
        }
        populate(board.D, K)
        return board
}

type Cell struct {
        D []*Cell
        Player int
        IsClaimed bool
}

type Rules struct {
        InARow int
}

type Direction []int
type Point []int
type Segment []Point

func (p Point) Clone() Point {
        point := make(Point, len(p))
        copy(point, p)
        return point
}

func NewSegment(p Point, d Direction, length int) Segment {
        segment := make(Segment, length)
        segment[0] = p.Clone()
        for i := 1; i < length; i++ {
                segment[i] = segment[i-1].Clone()
                for j, v := range d {
                        segment[i][j] += v
                }
        }
        return segment
}
