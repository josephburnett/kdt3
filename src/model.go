package kdt3

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
}
