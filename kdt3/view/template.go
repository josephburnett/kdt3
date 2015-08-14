package view

import (
        html "html/template"
)

var RootTemplate = html.Must(html.New("root").Parse(rootTemplateHTML))
const rootTemplateHTML = `
<html>
  <body>
    <h3>K-D Tic Tac Toe</h3>
    <p>Welcome {{.}}. Choose a game style below:</p>
    <ul>
      <li><a href="/game?handle=Player 1;playerCount=2;k=2;size=3;inarow=3">Classic (2D-2P)</a></li>
      <li><a href="/game?handle=Player 1;playerCount=2;k=3;size=4;inarow=4">Deep Thinker (3D-2P)</a></li>
      <li><a href="/game?handle=Player 1;playerCount=4;k=4;size=5;inarow=5">Preposterous Party (4D-4P)</a></li>
      <li><a href="/new">Custom</a></li>
    </ul>
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
{{ $inARow := .Rules.InARow }}
<html>
  <h3>New Game</h3>
  <body>
    <p>Share these links to play:</p>
    <ol>
    {{range .Players}}<li><a href="game/{{.PlayerId}}?message=Game on! Get {{ $inARow }} in a row to win.">{{.Handle}}</a></li>{{end}}
    </ol>
  </body>
</html>
`

var GetGameTemplate = html.Must(html.New("getgame").Parse(getGameTemplateHTML))
const getGameTemplateHTML = `
<html>
  <body>
    <h3>{{.PlayerHandle}}{{ if .MyTurn }} (your turn){{end}}</h3>
    {{if .Won}}<p>Game over!</p>{{else}}<p>{{.Message}}</p>{{end}}
    <div>{{.View}}</div>
    <div>{{.PlayerList}}</div>
  </body>
<html>
`
