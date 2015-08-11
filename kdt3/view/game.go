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
        tableBegin := "<table>"
        tableEnd := "</table>"
        rowBegin := "<tr>"
        rowEnd := "</tr>"
        columnBegin := "<td><div style=\"border: 1px solid;\">"
        columnEnd := "</div></td>"
        var recur func(*m.Cell, int) string
        recur = func(c *m.Cell, depth int) string {
                if depth == 0 {
                        if c.IsClaimed {
                                return strconv.Itoa(c.Player)
                        } else {
                                return "&nbsp"
                        }
                } else if depth % 2 == 0 {
                        table := tableBegin
                        for _, v := range c.D {
                                table += rowBegin + recur(v, depth-1) + rowEnd
                        }
                        table += tableEnd
                        return table
                } else {
                        columns := ""
                        for _, v := range c.D {
                                columns += columnBegin + recur(v, depth-1) + columnEnd
                        }
                        return columns
                }
        }
        if b.K %2 == 0 {
                return recur(b.D, b.K)
        } else {
                return tableBegin + "<tr>" + recur(b.D, b.K) + "</tr>" + tableEnd
        }
}
