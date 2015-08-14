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
        player := g.PlayerIndex()
        boardView := &ViewableBoard{g.Board, g.PlayerId, player}
        return template.HTML(boardView.View())
}

func (g *ViewableGame) PlayerIndex() int {
        for i, v := range g.PlayerIds {
                if g.PlayerId == v {
                        return i
                }
        }
        return g.Turn
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

func (g *ViewableGame) MyTurn() bool {
        return g.PlayerId == g.PlayerIds[g.Turn]
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
        Player int
}

func (b *ViewableBoard) View() string {
        tableBegin := "<table>"
        tableEnd := "</table>"
        rowBegin := "<tr>"
        rowEnd := "</tr>"
        columnBegin := "<td><div style=\"border: 1px solid;\">"
        columnEnd := "</div></td>"
        cellBegin := "<div style=\"width: 30px; height: 30px; text-align: center; line-height: 30px;\">"
        cellBeginWin := "<div style=\"width: 30px; height: 30px; text-align: center; line-height: 30px; background: #ABF095;\">"
        cellBeginLoss := "<div style=\"width: 30px; height: 30px; text-align: center; line-height: 30px; background: #F7C1C1;\">"
        cellBeginMyClaim := "<div style=\"width: 30px; height: 30px; text-align: center; line-height: 30px; background: #C1EAF7;\">"
        cellBeginYourClaim := "<div style=\"width: 30px; height: 30px; text-align: center; line-height: 30px; background: #E6E3E3;\">"
        cellEnd := "</div>"
        point := make(m.Point, b.K)
        var recur func(*m.Cell, int) string
        recur = func(c *m.Cell, depth int) string {
                if depth == 0 {
                        if c.IsClaimed {
                                if c.IsWon {
                                        if c.Player == b.Player {
                                                return cellBeginWin + strconv.Itoa(c.Player+1) + cellEnd
                                        } else {
                                                return cellBeginLoss + strconv.Itoa(c.Player+1) + cellEnd
                                        }
                                } else if c.Player == b.Player {
                                        return cellBeginMyClaim + strconv.Itoa(c.Player+1) + cellEnd
                                } else {
                                        return cellBeginYourClaim + strconv.Itoa(c.Player+1) + cellEnd
                                }
                        } else {
                                return "<a href=\"/move/" + b.PlayerId + "?point=" + point.String() +
                                       "\">" + cellBegin + cellEnd + "</a>"
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
