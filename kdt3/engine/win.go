package engine

import (
        "kdt3/model"
)

const (
        Decline = -1
        Neutral = 0
        Incline = 1
)

// for each cell, look each direction and count to in-a-row.

func eachDirection(K int, fn func ([]int)) {
        direction := make([]int, K)
        var recur func(fn func([]int), depth, dir int)
        recur = func(fn func([]int), depth, dir int) {
                direction[depth] = dir
                if depth == 0 {
                        fn(direction)
                } else {
                        recur(fn, depth-1, Decline)
                        recur(fn, depth-1, Neutral)
                        recur(fn, depth-1, Incline)
                }
        }
        recur(fn, K-1, Decline)
        recur(fn, K-1, Neutral)
        recur(fn, K-1, Incline)
}

func isWinningVector(K int, root *model.Cell, player int, rules *model.Rules, point, direction []int) bool {
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
