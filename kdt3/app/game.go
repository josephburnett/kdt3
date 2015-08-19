package kdt3

import (
        "net/http"
        "time"

        "appengine"
        "appengine/channel"
        "appengine/datastore"
        "appengine/memcache"
        "appengine/user"

        m "kdt3/model"
)

func getToken(c appengine.Context, playerId string) (string, error) {
        item, err := memcache.Get(c, "token::"+playerId)
        if err != nil {
                tok, err := channel.Create(c, playerId)
                if err != nil {
                        return "", err
                }
                memcache.Set(c, &memcache.Item{
                        Key: "token::"+playerId,
                        Value: []byte(tok),
                        Expiration: 2 * time.Hour,
                })
                return tok, nil
        }
        tok := string(item.Value)
        return tok, nil
}

func updateClients(c appengine.Context, game *m.Game) {
        for _, player := range game.Players {
                channel.SendJSON(c, player.PlayerId, "true")
        }
}

func loadGame(c appengine.Context, gameId, playerId string) (*m.Game, *m.Player, error) {
        gameKey := datastore.NewKey(c, "Game", gameId, 0, nil)
        game := &m.Game{}
        err := datastore.Get(c, gameKey, game)
        if err != nil {
                return nil, nil, err
        }
        if playerId == "" {
                return game, nil, nil
        }
        playerKey := datastore.NewKey(c, "Player", playerId, 0, gameKey)
        player := &m.Player{}
        err = datastore.Get(c, playerKey, player)
        if err != nil {
                return game, nil, err
        }
        if player.GameId != game.GameId {
                return game, nil, nil
        }
        return game, player, nil
}

func saveGame(c appengine.Context, game *m.Game) error {
        now, err := time.Now().MarshalText()
        if err != nil {
                return err
        }
        game.UpdatedDate = string(now)
        err = datastore.RunInTransaction(c, func(c appengine.Context) error {
                gameKey := datastore.NewKey(c, "Game", game.GameId, 0, nil)
                _, err := datastore.Put(c, gameKey, game)
                return err
        }, nil)
        return err
}

func createGame(c appengine.Context, r *http.Request) (*m.Game, error) {
        game, err := m.NewGame(r)
        if err != nil {
                return nil, err
        }
        u := user.Current(c)
        if u != nil {
                game.Creator = u.ID
        }
        err = datastore.RunInTransaction(c, func(c appengine.Context) error {
                gameKey := datastore.NewKey(c, "Game", game.GameId, 0, nil)
                playerKeys := make([]*datastore.Key, len(game.Players))
                for i, p := range game.Players {
                        playerKeys[i] = datastore.NewKey(c, "Player", p.PlayerId, 0, gameKey)
                }
                _, err := datastore.Put(c, gameKey, game)
                if err != nil {
                        return err
                }
                _, err = datastore.PutMulti(c, playerKeys, game.Players)
                return err
        }, nil)
        if err != nil {
                return nil, err
        }
        return game, nil
}
