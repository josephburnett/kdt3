package view

import (
        html "html/template"
)

var NewGameTemplate = html.Must(html.New("getroot").Parse(newGameTemplateHTML))
const newGameTemplateHTML = `
<html>
  <body>
    <h3>New Game</h3>
    <form action="/game" method="post">
      <div>Handle: <input type="text" name="handle"></div>
      <div>PlayerCount: <input type="text" name="playerCount"></div>
      <div>K: <input type="text" name="k"></div>
      <div>Size: <input type="text" name="size"></div>
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
    <p>Click <a href="/game/{{.GameId}}">{{.GameId}}</a> to play!</p>
  </body>
</html>
`

var GetGameTemplate = html.Must(html.New("getgame").Parse(getGameTemplateHTML))
const getGameTemplateHTML = `
<html>
  <body>
    <h3>{{.GameId}} ({{.Turn}})</h3>
    <div>{{.View}}</div>
  </body>
<html>
`
