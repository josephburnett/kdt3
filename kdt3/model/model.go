package model

type Player struct {
        PlayerId string
        GameId string
        Handle string
}

type Game struct {
        GameId string
        PlayerIds []string
        Owner int
        Turn int
        Win [][]int
        Board *Board
}

type Board struct {
        K int
        D *Cell
        Size int
}

type Cell struct {
        D []*Cell
        Player int
        IsClaimed bool
}

type Rules struct {
        InARow int
}
