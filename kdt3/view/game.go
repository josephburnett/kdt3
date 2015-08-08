package view

import (
        "strconv"
        "html/template"

        m "kdt3/model"
)

type ViewableGame struct {
        *m.Game
        BoardHTML template.HTML
}

type ViewableBoard struct {
        *m.Board
}

func (b *ViewableBoard) View() string {
        var recur func(*m.Cell, int) string
        recur = func(c *m.Cell, depth int) string {
                if depth == 0 {
                        if c.IsClaimed {
                                return strconv.Itoa(c.Player)
                        } else {
                                return "_"
                        }
                } else {
                        table := "<table><tr>"
                        for _, v := range c.D {
                                table += "<td>" + recur(v, depth-1) + "</td>"
                        }
                        table += "</tr></table>"
                        return table
                }
        }
        return recur(b.D, b.K)
}
