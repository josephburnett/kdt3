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
                        players += "<li>" + p.Handle + "</li>"
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
        point := make(m.Point, b.K)
        var recur func(*m.Cell, int) string
        recur = func(c *m.Cell, depth int) string {
                if depth == 0 {
                        classes := "cell"
                        if c.IsClaimed {
                                if c.IsWon {
                                        if c.Player == b.Player {
                                                classes += " win"
                                        } else {
                                                classes += " loss"
                                        }
                                } else if c.Player == b.Player {
                                        classes += " mine"
                                } else {
                                        classes += " yours"
                                }
                                return "<div class=\"" + classes + "\">" + strconv.Itoa(c.Player+1) + "</div>"
                        } else {
                                return "<a href=\"/move/" + b.PlayerId + "?point=" + point.String() +
                                       "\"><div class=\"" + classes + "\"></div></a>"
                        }
                } else if depth % 2 == 0 {
                        table := "<table>"
                        for i, v := range c.D {
                                point[b.K-depth] = i
                                table += "<tr>" + recur(v, depth-1) + "</tr>"
                        }
                        table += "</table>"
                        return table
                } else {
                        columns := ""
                        for i, v := range c.D {
                                point[b.K-depth] = i
                                columns += "<td><div class=\"col\">" + recur(v, depth-1) + "</div></td>"
                        }
                        return columns
                }
        }
        if b.K %2 == 0 {
                return recur(b.D, b.K)
        } else {
                return "<table>" + "<tr>" + recur(b.D, b.K) + "</tr>" + "</table>"
        }
}
