package kdt3

import (
        "net/http"

        "appengine"
        "appengine/user"

        "kdt3/engine"
        m "kdt3/model"
        "kdt3/view"
)

func init() {
        http.HandleFunc("/", getRoot)
        http.HandleFunc("/new", getNew)
        http.HandleFunc("/game", postGame)
        http.HandleFunc("/game/", getGame)
        http.HandleFunc("/move/", postMove)
}

func getRoot(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        u := user.Current(c)
        if u == nil {
                redirectLogin(c, w, r)
                return
        }
        err := view.RootTemplate.Execute(w, u)
        if internalError(w, err) {
                return
        }
}

func getNew(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        u := user.Current(c)
        if u == nil {
                redirectLogin(c, w, r)
                return
        }
        message := r.FormValue("message")
        err := view.NewGameTemplate.Execute(w, message)
        if internalError(w, err) {
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
                http.Redirect(w, r, "/new?message="+err.Error(), http.StatusFound)
                return
        }
        err = view.PostGameTemplate.Execute(w, game)
        if internalError(w, err) {
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
        if internalError(w, err) {
                return
        }
        gameView := view.NewViewableGame(game, id)
        gameView.Message = r.FormValue("message")
        err = view.GetGameTemplate.Execute(w, gameView)
        if internalError(w, err) {
                return
        }
}

func postMove(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)
        u := user.Current(c)
        if u == nil {
                redirectLogin(c, w, r)
                return
        }
        id := r.URL.Path[len("/move/"):]
        game, err := loadGame(c, id)
        if internalError(w, err) {
                return
        }
        if game.Won {
                http.Redirect(w, r, "/game/"+id, http.StatusFound)
                return
        }
        point, err := m.ParsePoint(game.Board.K, game.Board.Size, r.FormValue("point"))
        if internalError(w, err) {
                return
        }
        gameMove := &engine.MovableGame{game}
        err = gameMove.Move(id, point)
        if err != nil {
                http.Redirect(w, r, "/game/"+id+"?message="+err.Error(), http.StatusFound)
                return
        }
        gameWin := &engine.WinnableGame{game}
        if gameWin.IsWin() {
                game.Won = true
        } else {
                gameMove.AdvanceTurn()
        }
        err = saveGame(c, gameMove.Game)
        if internalError(w, err) {
                return
        }
        http.Redirect(w, r, "/game/"+id+"?message=Move accepted.", http.StatusFound)
}

func redirectLogin(c appengine.Context, w http.ResponseWriter, r *http.Request) {
        url, err := user.LoginURL(c, r.URL.String())
        if internalError(w, err) {
                return
        }
        w.Header().Set("location", url)
        w.WriteHeader(http.StatusFound)
}

func internalError(w http.ResponseWriter, err error) bool {
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return true
        } else {
                return false
        }
}
