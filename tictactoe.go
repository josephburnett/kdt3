package tictactoe

import (
	"fmt"
	"crypto/rand"
	"encoding/base32"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"appengine"
	"appengine/memcache"
	"appengine/user"
)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/new", newGame)
	http.HandleFunc("/game", postGame)
	http.HandleFunc("/game/", getGame)
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		redirectLogin(c, w, r)
		return
	}
	fmt.Fprintf(w, "Hello %v.", u)
}

func newGame(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		redirectLogin(c, w, r)
		return
	}
	fmt.Fprint(w, newGameForm)
}

const newGameForm = `
<html>
  <body>
    <form action="/game" method="post">
      <div>Handle: <input type="text" name="handle"></div>
      <div>PlayerCount: <input type="text" name="playerCount"></div>
      <div><input type="submit" value="Create"></div>
    </form>
  </body>
</html>
`

func postGame(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		redirectLogin(c, w, r)
		return
	}
	game, err := createGame(c, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = postGameTemplate.Execute(w, game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

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

var postGameTemplate = template.Must(template.New("postgame").Parse(postGameTemplateHTML))

const postGameTemplateHTML = `
<html>
  <body>
    <p>Click <a href="/game/{{.GameId}}">{{.GameId}}</a> to play!</p>
  </body>
</html>
`

func getGame (w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		redirectLogin(c, w, r)
		return
	}
	id := r.URL.Path[len("/game/"):]
	game := &Game{}
	_, err := memcache.JSON.Get(c, id, game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = getGameTemplate.Execute(w, game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var getGameTemplate = template.Must(template.New("getgame").Parse(getGameTemplateHTML))

const getGameTemplateHTML = `
<html>
  <body>
    <pre>{{.}}</pre>
  </body>
<html>
`

func newId() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	str := base32.StdEncoding.EncodeToString(bytes)
	return strings.TrimSuffix(str, "======")
}

func redirectLogin(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	url, err := user.LoginURL(c, r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("location", url)
	w.WriteHeader(http.StatusFound)
	return
}

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
