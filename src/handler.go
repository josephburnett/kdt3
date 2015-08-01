package kdt3

import (
        "fmt"
        "net/http"

        "appengine"
        "appengine/memcache"
        "appengine/user"
)

func init() {
        http.HandleFunc("/", getRoot)
        http.HandleFunc("/new", getNew)
        http.HandleFunc("/game", postGame)
        http.HandleFunc("/game/", getGame)
}

func getRoot(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        u := user.Current(c)
        if u == nil {
                redirectLogin(c, w, r)
                return
        }
        fmt.Fprintf(w, "Hello %v.", u)
}

func getNew(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        u := user.Current(c)
        if u == nil {
                redirectLogin(c, w, r)
                return
        }
        fmt.Fprint(w, newGameForm)
}

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

func getGame(w http.ResponseWriter, r *http.Request) {
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
