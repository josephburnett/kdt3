package engine

import (
        "errors"

        m "kdt3/model"
)

type MovableGame struct {
        *m.Game
}

func (g *MovableGame) Move(playerId string, point m.Point) error {
        if playerId != g.PlayerIds[g.Turn] {
                return errors.New("Invalid move: out of turn.")
        }
        cell := g.Board.D
        for _, v := range point {
                cell = cell.D[v]
        }
        if cell.IsClaimed {
                return errors.New("Invalid move: cell already claimed.")
        }
        cell.IsClaimed = true
        cell.Player = g.Turn
        return nil
}

func (g *MovableGame) AdvanceTurn() {
        g.Turn = (g.Turn + 1) % len(g.PlayerIds)
}
