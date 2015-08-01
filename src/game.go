package kdt3

import (
        "net/http"
        "strconv"
        "time"

        "appengine"
        "appengine/memcache"
)

func createGame(c appengine.Context, r *http.Request) (*Game, error) {
        // parameters
        playerCount, err := strconv.Atoi(r.FormValue("playerCount"))
        if err != nil {
                return nil, err
        }
        playerIds := make([]string, playerCount)
        handles := make([]string, playerCount)
        for i, _ := range playerIds {
                playerIds[i] = newId()
                handles[i] = "Player " + strconv.Itoa(i+1)
        }
        handles[0] = r.FormValue("handle")
        // objects
        items := make([]*memcache.Item, playerCount+1)
        game := &Game {
                GameId: newId(),
                PlayerIds: playerIds,
                Owner: 0,
                Turn: 0,
        }
        items[playerCount] = &memcache.Item {
                Key: game.GameId,
                Object: game,
                Expiration: 25 * time.Hour,
        }
        for i, _ := range playerIds {
                player := &Player {
                        PlayerId: newId(),
                        GameId: game.GameId,
                        Handle: handles[i],
                }
                items[i] = &memcache.Item {
                        Key: player.PlayerId,
                        Object: player,
                        Expiration: 24 * time.Hour,
                }
        }
        // persist
        err = memcache.JSON.AddMulti(c, items)
        if err != nil {
                return nil, err
        }
        return game, err
}
