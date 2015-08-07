package kdt3

import (
        html "html/template"
)

const newGameForm = `
<html>
  <body>
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

var postGameTemplate = html.Must(html.New("postgame").Parse(postGameTemplateHTML))
const postGameTemplateHTML = `
<html>
  <body>
    <p>Click <a href="/game/{{.GameId}}">{{.GameId}}</a> to play!</p>
  </body>
</html>
`

var getGameTemplate = html.Must(html.New("getgame").Parse(getGameTemplateHTML))
const getGameTemplateHTML = `
<html>
  <body>
    <pre>{{.}}</pre>
  </body>
<html>
`
