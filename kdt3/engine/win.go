package engine

import (
        "kdt3/model"
)

const (
        Decline = -1
        Neutral = 0
        Incline = 1
)

type WinnableGame struct {
        *model.Game
}

func (g *WinnableGame) IsWin() bool {
        boardWin := &WinnableBoard{g.Board}
        wins := boardWin.GetWins(g.TurnOrder, g.Rules)
        if len(wins) == 0 {
                return false
        }
        boardWin.AnnotateWins(wins)
        return true
}

type WinnableBoard struct {
        *model.Board
}

func (b *WinnableBoard) AnnotateWins(wins []model.Segment) {
        for _, segment := range wins {
                for _, point := range segment {
                        cell := b.D
                        for _, i := range point {
                                cell = cell.D[i]
                        }
                        cell.IsWon = true
                }
        }
}

func (b *WinnableBoard) GetWins(player int, rules *model.Rules) (wins []model.Segment) {
        EachPoint(b.K, b.D, func(p model.Point) {
                EachDirection(b.K, func(d model.Direction) {
                        if IsWinningVector(b.K, b.D, player, rules, p, d) {
                                wins = append(wins, model.NewSegment(p, d, rules.InARow))
                        }
                })
        })
        return
}

func EachPoint(K int, root *model.Cell, fn func(model.Point)) {
        point := make(model.Point, K)
        var recur func(*model.Cell, int)
        recur = func(c *model.Cell, depth int) {
                if depth == 0 {
                        fn(point)
                } else {
                        for i, d := range c.D {
                                point[depth-1] = i
                                recur(d, depth-1)
                        }
                }
        }
        recur(root, K)
}

func EachDirection(K int, fn func (model.Direction)) {
        direction := make(model.Direction, K)
        var recur func(int, int)
        recur = func(depth, dir int) {
                direction[depth] = dir
                if depth == 0 {
                        fn(direction)
                } else {
                        recur(depth-1, Decline)
                        recur(depth-1, Neutral)
                        recur(depth-1, Incline)
                }
        }
        recur(K-1, Decline)
        recur(K-1, Neutral)
        recur(K-1, Incline)
}

func IsWinningVector(K int, root *model.Cell, player int, rules *model.Rules, point, direction []int) bool {
        entirelyNeutral := true
        for _, v := range direction {
                if v != Neutral {
                        entirelyNeutral = false
                }
        }
        if entirelyNeutral {
                return false
        }
        nextPoint := make([]int, K)
        copy(nextPoint, point)
        for i := 0; i < rules.InARow; i++ {
                // check out of bounds
                for _, v := range nextPoint {
                        if v < 0 || v > len(root.D)-1 {
                                return false
                        }
                }
                // check claim
                cell := root
                for _, j := range nextPoint {
                        cell = cell.D[j]
                }
                if !cell.IsClaimed || cell.Player != player {
                        return false
                }
                // find next cell
                for j, _ := range nextPoint {
                        nextPoint[j] += direction[j]
                }
        }
        return true
}
