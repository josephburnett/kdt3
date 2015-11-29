package ai

import (
        "math"

        m "kdt3/model"
        e "kdt3/engine"
)

// Count the number of open vectors for player minus all other open vectors
// for other players.  If any opponent has a winning vector, score is MinInt.
// Otherwise if plaher has winning vector, score is MaxInt.
func Evaluate(board *m.Board, rules *m.Rules, playerCount, player int) (score int) {
        win := false
        loss := false
        for i := 0; i < playerCount; i++ {
                open, wins := countOpenVectors(board, i, rules)
                if i == player {
                        score += open
                        if wins > 0 {
                                win = true
                        }
                } else {
                        score -= open
                        if wins > 0 {
                                loss = true
                        }
                }
        }
        if win {
                score = math.MaxInt64
        }
        if loss {
                score = math.MinInt64
        }
        return
}

// Count the number of open and winning vectors for a given player.
func countOpenVectors(board *m.Board, playerOrder int, rules *m.Rules) (open, wins int) {
        e.EachPoint(board.K, board.D, func(p m.Point) {
                e.EachDirection(board.K, func(d m.Direction) {
                        if isOpenVector(board.K, board.D, playerOrder, rules, p, d) {
                                open += 1
                        }
                        if e.IsWinningVector(board.K, board.D, playerOrder, rules, p, d) {
                                wins += 1
                        }
                })
        })
        return
}

func isOpenVector(K int, root *m.Cell, player int, rules *m.Rules, point, direction []int) bool {
        entirelyNeutral := true
        for _, v := range direction {
                if v != e.Neutral {
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
                if cell.IsClaimed && cell.Player != player {
                        return false
                }
                // find next cell
                for j, _ := range nextPoint {
                        nextPoint[j] += direction[j]
                }
        }
        return true
}
