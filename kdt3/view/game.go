package view

import (
        "strconv"
        "html/template"

        m "kdt3/model"
)

type ViewableGame struct {
        *m.Game
        Viewer *m.Player
        IsMyTurn bool
        Message string
}

func NewViewableGame(game *m.Game, playerId string) *ViewableGame {
        viewableGame := &ViewableGame{
                Game: game,
        }
        for _, v := range game.Players {
                if playerId == v.PlayerId {
                        viewableGame.Viewer = v
                }
        }
        if playerId == game.TurnId {
                viewableGame.IsMyTurn = true
        }
        return viewableGame
}

func (g *ViewableGame) View() template.HTML {
        boardView := &ViewableBoard{g.Board, g.Viewer}
        return template.HTML(boardView.View())
}

func (g *ViewableGame) PlayerList() template.HTML {
        players := "<ol>"
        for i, p := range g.Players {
                if i == g.TurnOrder {
                        players += "<li><b>" + p.Handle + "</b></li>"
                } else {
                        players += "<li>" + p.Handle + "</li>"
                }
        }
        players += "</ol>"
        return template.HTML(players)
}

type ViewableBoard struct {
        *m.Board
        Viewer *m.Player
}

func (b *ViewableBoard) View() string {
        point := make(m.Point, b.K)
        var recur func(*m.Cell, int) string
        recur = func(c *m.Cell, depth int) string {
                if depth == 0 {
                        classes := "cell"
                        if c.IsClaimed {
                                if c.IsWon {
                                        if c.Player == b.Viewer.PlayerOrder {
                                                classes += " win"
                                        } else {
                                                classes += " loss"
                                        }
                                } else if c.Player == b.Viewer.PlayerOrder {
                                        classes += " mine"
                                } else {
                                        classes += " yours"
                                }
                                return "<div class=\"" + classes + "\">" + strconv.Itoa(c.Player+1) + "</div>"
                        } else {
                                return "<a href=\"/move/" + b.Viewer.PlayerId + "?point=" + point.String() +
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
