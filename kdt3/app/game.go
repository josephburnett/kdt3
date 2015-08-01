package kdt3

import (
        "net/http"
        "strconv"
        "time"

        "appengine"
        "appengine/memcache"

        "kdt3/model"
)

func createGame(c appengine.Context, r *http.Request) (*model.Game, error) {
        // parameters
        playerCount, err := strconv.Atoi(r.FormValue("playerCount"))
        if err != nil {
                return nil, err
        }
        gameId := newId()
        playerIds := make([]string, playerCount)
        handles := make([]string, playerCount)
        for i, _ := range playerIds {
                playerIds[i] = newId()
                handles[i] = "Player " + strconv.Itoa(i+1)
        }
        handles[0] = r.FormValue("handle")
        // objects
        items := make([]*memcache.Item, playerCount+1)
        for i, _ := range playerIds {
                player := &model.Player {
                        PlayerId: playerIds[i],
                        GameId: gameId,
                        Handle: handles[i],
                }
                items[i] = &memcache.Item {
                        Key: player.PlayerId,
                        Object: player,
                        Expiration: 24 * time.Hour,
                }
        }
        game := &model.Game {
                GameId: gameId,
                PlayerIds: playerIds,
                Owner: 0,
                Turn: 0,
        }
        items[playerCount] = &memcache.Item {
                Key: gameId,
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