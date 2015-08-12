package engine

import (
        m "kdt3/model"
)

type MovableGame struct {
        *m.Game
}

func (g *MovableGame) Move(m.Point) error {
        return nil
}
