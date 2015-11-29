package ai

import (
        m "kdt3/model"
)

type Player interface {
        Name() string
        Play(*m.Game) *m.Point
}
