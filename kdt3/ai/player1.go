package ai

import (
        m "kdt3/model"
)

type Player1 struct {}

func (p *Player1) Name() string {
        return "Minimax AI Player"
}

func (p *Player1) Play(*m.Game) *m.Point {
        return nil
}
