package kdt3

import (
        "errors"
        "net/http"
        "strconv"
        "time"

        "appengine"
        "appengine/memcache"

        m "kdt3/model"
)

func loadGame(c appengine.Context, playerId string) (*m.Game, error) {
        player := &m.Player{}
        _, err := memcache.JSON.Get(c, "player::" + playerId, player)
        if err != nil {
                return nil, err
        }
        game := &m.Game{}
        _, err = memcache.JSON.Get(c, "game::" + player.GameId, game)
        if err != nil {
                return nil, err
        }
        return game, nil
}

func saveGame(c appengine.Context, game *m.Game) error {
        item := &memcache.Item {
                Key: "game::" + game.GameId,
                Object: game,
                Expiration: 25 * time.Hour,
        }
        err := memcache.JSON.Set(c, item)
        return err
}

func createGame(c appengine.Context, r *http.Request) (*m.Game, error) {
        // parameters
        playerCount, err := strconv.Atoi(r.FormValue("playerCount"))
        if err != nil {
                return nil, err
        }
        if playerCount < 2 || playerCount > 10 {
                return nil, errors.New("Player Count must be between 2 and 10.")
        }
        K, err := strconv.Atoi(r.FormValue("k"))
        if err != nil {
                return nil, err
        }
        if K < 2 || K > 5 {
                return nil, errors.New("K must be between 2 and 5")
        }
        size, err := strconv.Atoi(r.FormValue("size"))
        if err != nil {
                return nil, err
        }
        if size < 2 || size > 5 {
                return nil, errors.New("Size must be between 2 and 5")
        }
        inARow, err := strconv.Atoi(r.FormValue("inarow"))
        if err != nil {
                return nil, err
        }
        if inARow < 2 || inARow > size {
                return nil, errors.New("In a row must be between 2 and " + strconv.Itoa(size))
        }
        gameId := newId()
        playerIds := make([]string, playerCount)
        players := make([]*m.Player, playerCount)
        handles := make([]string, playerCount)
        for i, _ := range playerIds {
                playerIds[i] = newId()
                handles[i] = "Player " + strconv.Itoa(i+1)
        }
        handles[0] = r.FormValue("handle")
        // objects
        items := make([]*memcache.Item, playerCount+1)
        for i, _ := range playerIds {
                player := &m.Player {
                        PlayerId: playerIds[i],
                        GameId: gameId,
                        Handle: handles[i],
                }
                players[i] = player
                items[i] = &memcache.Item {
                        Key: "player::" + player.PlayerId,
                        Object: player,
                        Expiration: 24 * time.Hour,
                }
        }
        board := m.NewBoard(K, size)
        game := &m.Game {
                GameId: gameId,
                PlayerIds: playerIds,
                Players: players,
                Owner: 0,
                Turn: 0,
                Board: board,
                Rules: &m.Rules{InARow: inARow},
        }
        items[playerCount] = &memcache.Item {
                Key: "game::" + gameId,
                Object: game,
                Expiration: 25 * time.Hour,
        }
        // persist
        err = memcache.JSON.AddMulti(c, items)
        if err != nil {
                return nil, err
        }
        return game, err
}
