package view

import (
        "strconv"
        "html/template"

        m "kdt3/model"
)

type ViewableGame struct {
        *m.Game
        PlayerId string
        Message string
}

func (g *ViewableGame) View() template.HTML {
        boardView := &ViewableBoard{g.Board, g.PlayerId}
        return template.HTML(boardView.View())
}

func (g *ViewableGame) PlayerList() template.HTML {
        players := "<ol>"
        for i, p := range g.Players {
                if i == g.Turn {
                        players += "<li><b>" + p.Handle + "</b></li>"
                } else {
                        players += "<li>" + p.Handle + "</b></li>"
                }
        }
        players += "</ol>"
        return template.HTML(players)
}

func (g *ViewableGame) PlayerHandle() string {
        for i, p := range g.PlayerIds {
                if p == g.PlayerId {
                        return g.Players[i].Handle
                }
        }
        return "Unknown"
}

type ViewableBoard struct {
        *m.Board
        PlayerId string
}

func (b *ViewableBoard) View() string {
        tableBegin := "<table>"
        tableEnd := "</table>"
        rowBegin := "<tr>"
        rowEnd := "</tr>"
        columnBegin := "<td><div style=\"border: 1px solid;\">"
        columnEnd := "</div></td>"
        cellBegin := "<div style=\"width: 30px; height: 30px;\">"
        cellEnd := "</div>"
        point := make(m.Point, b.K)
        var recur func(*m.Cell, int) string
        recur = func(c *m.Cell, depth int) string {
                if depth == 0 {
                        if c.IsClaimed {
                                return cellBegin + strconv.Itoa(c.Player+1) + cellEnd
                        } else {
                                return cellBegin + "<a href=\"/move/" +
                                       b.PlayerId + "?point=" + point.String() +
                                       "\">&nbsp</a>" + cellEnd
                        }
                } else if depth % 2 == 0 {
                        table := tableBegin
                        for i, v := range c.D {
                                point[b.K-depth] = i
                                table += rowBegin + recur(v, depth-1) + rowEnd
                        }
                        table += tableEnd
                        return table
                } else {
                        columns := ""
                        for i, v := range c.D {
                                point[b.K-depth] = i
                                columns += columnBegin + recur(v, depth-1) + columnEnd
                        }
                        return columns
                }
        }
        if b.K %2 == 0 {
                return recur(b.D, b.K)
        } else {
                return tableBegin + rowBegin + recur(b.D, b.K) + rowEnd + tableEnd
        }
}
