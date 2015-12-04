package ai

import (
        e "kdt3/engine"
        m "kdt3/model"
)

func GenerateMoves(board *m.Board) (moves []*m.Point) {
        e.EachPoint(board.K, board.D, func(p m.Point) {
                cell := board.D
                for _, i := range p {
                        cell = cell.D[i]
                }
                if !cell.IsClaimed {
                        m := p.Clone()
                        moves = append(moves, &m)
                }
        })
        return
}
