package model

import (
        "errors"
        "net/http"
        "strconv"
        "strings"
)

type Player struct {
        PlayerOrder int
        PlayerId string
        GameId string
        Handle string
}

type Game struct {
        GameId string
        Players []*Player
        TurnOrder int
        TurnId string
        Won bool
        Board *Board
        Rules *Rules
}

func NewGame(r *http.Request) (*Game, error) {
        playerCount, err := strconv.Atoi(r.FormValue("playerCount"))
        if err != nil {
                return nil, err
        }
        if playerCount < 2 || playerCount > 10 {
                return nil, errors.New("Player Count must be between 2 and 10.")
        }
        K, err := strconv.Atoi(r.FormValue("k"))
        if err != nil {
                return nil, err
        }
        if K < 2 || K > 5 {
                return nil, errors.New("K must be between 2 and 5")
        }
        size, err := strconv.Atoi(r.FormValue("size"))
        if err != nil {
                return nil, err
        }
        if size < 2 || size > 5 {
                return nil, errors.New("Size must be between 2 and 5")
        }
        inARow, err := strconv.Atoi(r.FormValue("inarow"))
        if err != nil {
                return nil, err
        }
        if inARow < 2 || inARow > size {
                return nil, errors.New("In a row must be between 2 and " + strconv.Itoa(size))
        }
        gameId := newId()
        players := make([]*Player, playerCount)
        for i, _ := range players {
                players[i] = &Player{
                        PlayerId: newId(),
                        GameId: gameId,
                        PlayerOrder: i,
                        Handle: "Player " + strconv.Itoa(i+1),
                }
        }
        game := &Game{
                GameId: gameId,
                Players: players,
                TurnOrder: 0,
                TurnId: players[0].PlayerId,
                Won: false,
                Board: NewBoard(K, size),
                Rules: &Rules{InARow: inARow},
        }
        return game, nil
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
        IsWon bool
}

type Rules struct {
        InARow int
}

type Direction []int
type Point []int
type Segment []Point

func ParsePoint(K, size int, s string) (Point, error) {
        ss := strings.Split(s, ",")
        if len(ss) != K {
                return nil, errors.New("Invalid point: incorrect number of dimensions.")
        }
        p := make(Point, K)
        for i, v := range ss {
                d, err := strconv.Atoi(v)
                if err != nil {
                        return nil, err
                }
                if d < 0 || d > size-1 {
                        return nil, errors.New("Invalid point: out of bounds.")
                }
                p[i] = d
        }
        return p, nil
}

func (p Point) Clone() Point {
        point := make(Point, len(p))
        copy(point, p)
        return point
}

func (p Point) String() string {
        ss := make([]string, len(p))
        for i, v := range p {
                ss[i] = strconv.Itoa(v)
        }
        return strings.Join(ss,",")
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
