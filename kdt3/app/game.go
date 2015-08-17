package kdt3

import (
        "net/http"
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
        game, err := m.NewGame(r)
        if err != nil {
                return nil, err
        }
        items := make([]*memcache.Item, len(game.Players)+1)
        for i, p := range game.Players {
                items[i] = &memcache.Item {
                        Key: "player::" + p.PlayerId,
                        Object: p,
                        Expiration: 24 * time.Hour,
                }
        }
        items[len(game.Players)] = &memcache.Item {
                Key: "game::" + game.GameId,
                Object: game,
                Expiration: 25 * time.Hour,
        }
        err = memcache.JSON.AddMulti(c, items)
        if err != nil {
                return nil, err
        }
        return game, err
}
