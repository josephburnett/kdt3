package kdt3

import (
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
        K, err := strconv.Atoi(r.FormValue("k"))
        if err != nil {
                return nil, err
        }
        size, err := strconv.Atoi(r.FormValue("size"))
        if err != nil {
                return nil, err
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
