package view

import (
        html "html/template"
)

var RootTemplate = html.Must(html.New("root").Parse(rootTemplateHTML))
const rootTemplateHTML = `
<html>
  <body>
    <h3>K-D Tic Tac Toe</h3>
    <p>Welcome {{.}}.  Click <a href="/new">here</a> to play!</p>
  </body>
</html>
`

var NewGameTemplate = html.Must(html.New("getnew").Parse(newGameTemplateHTML))
const newGameTemplateHTML = `
<html>
  <body>
    <h3>New Game</h3>
    <p>{{.}}</p>
    <form action="/game" method="post">
      <div>Handle: <input type="text" name="handle"></div>
      <div>PlayerCount: <input type="text" name="playerCount"></div>
      <div>K: <input type="text" name="k"></div>
      <div>Size: <input type="text" name="size"></div>
      <div>In a row: <input type="text" name="inarow"></div>
      <div><input type="submit" value="Create"></div>
    </form>
  </body>
</html>
`

var PostGameTemplate = html.Must(html.New("postgame").Parse(postGameTemplateHTML))
const postGameTemplateHTML = `
<html>
  <h3>New Game</h3>
  <body>
    <p>Share these links to play:</p>
    <ol>
    {{range .Players}}<li><a href="game/{{.PlayerId}}">{{.Handle}}</a></li>{{end}}
    </ol>
  </body>
</html>
`

var GetGameTemplate = html.Must(html.New("getgame").Parse(getGameTemplateHTML))
const getGameTemplateHTML = `
<html>
  <body>
    <h3>{{.PlayerHandle}}</h3>
    <p>{{.Message}}</p>
    {{if .Won}}<p>Game over!</p>{{end}}
    <div>{{.View}}</div>
    <div>{{.PlayerList}}</div>
  </body>
<html>
`
