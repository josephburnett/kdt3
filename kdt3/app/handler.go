package kdt3

import (
        "fmt"
        "net/http"

        "appengine"
        "appengine/user"

        "kdt3/view"
)

func init() {
        http.HandleFunc("/", getRoot)
        http.HandleFunc("/new", getNew)
        http.HandleFunc("/game", postGame)
        http.HandleFunc("/game/", getGame)
        //http.HandleFunc("/move/", postMove)
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
        err := view.NewGameTemplate.Execute(w, nil)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
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
                return
        }
        err = view.PostGameTemplate.Execute(w, game)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
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
        game, err := loadGame(c, id)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        gameView := &view.ViewableGame{game}
        err = view.GetGameTemplate.Execute(w, gameView)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
}
/*
func postMove(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        u := user.Current(c)
        if u == nil {
                redirectLogin(c, w, r)
                return
        }
        id := r.URL.Path[len("/move/"):]
        game := &model.Game{}
        _, err := memcache.JSON.Get(c, id, game)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        gameMove := &engine.MovableGame{game}
        //
}
*/
func redirectLogin(c appengine.Context, w http.ResponseWriter, r *http.Request) {
        url, err := user.LoginURL(c, r.URL.String())
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        w.Header().Set("location", url)
        w.WriteHeader(http.StatusFound)
}
