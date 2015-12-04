package ai

import (
        "testing"

        m "kdt3/model"
)

type PointList []*m.Point
func (a PointList) setEquals(b PointList) bool {
        if len(a) != len(b) {
                return false
        }
        for _,p1 := range a {
                has := false
                for _,p2 := range b {
                        if len(*p1) == len(*p2) {
                                for i,d := range *p1 {
                                        if d == (*p2)[i] {
                                                has = true
                                        }
                                }
                        }
                }
                if has == false {
                        return false
                }
        }
        return true
}

func TestPointListSetEquals(t *testing.T) {
        a := PointList([]*m.Point{
                &m.Point{1,0},
        })
        b := PointList([]*m.Point{
                &m.Point{0,1},
        })
        c := PointList([]*m.Point{
                &m.Point{1,0},
                &m.Point{0,1},
        })
        d := PointList([]*m.Point{
                &m.Point{0,1},
                &m.Point{1,0},
        })
        if true == a.setEquals(b) {
                t.Errorf("a.setEquals(b) is true")
        }
        if true == a.setEquals(c) {
                t.Errorf("a.setEquals(c) is true")
        }
        if false == c.setEquals(d) {
                t.Errorf("c.setEquals(d) is false")
        }
        if false == d.setEquals(c) {
                t.Errorf("d.setEquals(c) is false")
        }
}

func TestGenerateMovesK3D3(t *testing.T) {
        K := 3
        D := 3
        board := m.NewBoard(K, D)
        expected := PointList([]*m.Point{
                &m.Point{0,0,0},
                &m.Point{0,0,1},
                &m.Point{0,0,2},
                &m.Point{0,1,0},
                &m.Point{0,1,1},
                &m.Point{0,1,2},
                &m.Point{0,2,0},
                &m.Point{0,2,1},
                &m.Point{0,2,2},
                &m.Point{1,0,0},
                &m.Point{1,0,1},
                &m.Point{1,0,2},
                &m.Point{1,1,0},
                &m.Point{1,1,1},
                &m.Point{1,1,2},
                &m.Point{1,2,0},
                &m.Point{1,2,1},
                &m.Point{1,2,2},
                &m.Point{2,0,0},
                &m.Point{2,0,1},
                &m.Point{2,0,2},
                &m.Point{2,1,0},
                &m.Point{2,1,1},
                &m.Point{2,1,2},
                &m.Point{2,2,0},
                &m.Point{2,2,1},
                &m.Point{2,2,2},
        })
        moves := PointList(GenerateMoves(board))
        if !moves.setEquals(expected) {
                t.Errorf("Unexpected moves: %v", moves)
        }
}

func TestGenerateMovesK2D3(t *testing.T) {
        K := 2
        D := 3
        board := m.NewBoard(K, D)
        board.D.D[0].D[0].IsClaimed = true
        board.D.D[1].D[2].IsClaimed = true
        expected := PointList([]*m.Point{
                &m.Point{0,1},
                &m.Point{0,2},
                &m.Point{1,0},
                &m.Point{1,1},
                &m.Point{2,0},
                &m.Point{2,1},
                &m.Point{2,2},
        })
        moves := PointList(GenerateMoves(board))
        if !moves.setEquals(expected) {
                t.Errorf("Unexpected moves: %v", moves)
        }
}
