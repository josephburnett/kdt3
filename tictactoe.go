package tictactoe

import (
	"fmt"
	"crypto/rand"
	"encoding/base32"
	"html/template"
	"net/http"
	"strings"

	"appengine"
	"appengine/user"
)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/new", newGame)
	http.HandleFunc("/game", postGame)
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
	game := &Game {
		Creator: r.FormValue("handle"),
		Id: newId(),
	}
	err := postGameTemplate.Execute(w, game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var postGameTemplate = template.Must(template.New("postgame").Parse(postGameTemplateHTML))

const postGameTemplateHTML = `
<html>
  <body>
    <p>This game was started by {{.Creator}} and the id is {{.Id}}</p>
  </body>
</html>
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

type Game struct {
	Creator string
	Id string
}
