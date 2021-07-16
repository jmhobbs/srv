package main

import "html/template"

var Template *template.Template


func init () {
  Template = template.Must(template.New("index").Parse(`<!doctype html>
<html lang="en">
  <head>
    <title>{{.Path}}</title>
    <style>
    html, body {
      margin: 0;
      padding: 0;
    }
    h1 {
      width: 100%;
      background: #CCC;
      padding: 10px 20px;
      margin: 0;
    }
    div {
      display: block;
      border-bottom: 1px solid #DDD;
      padding: 5px 0 5px 45px;
      margin: 5px 0;
      font-size: 1.25em;
      background-repeat: no-repeat;
      background-position: 10px center;
      clear: both;
    }
    div a {
      text-decoration: none;
      color: blue;
    }
    div.parent {
      border-color: black;
    }
    div span {
      float: right;
      margin-right: 10px;
    }
    a:hover {
      color: purple;
      text-decoration: underline;
    }
    /*
      Icons by denimao from the Noun Project
      https://thenounproject.com/denimao/collection/user-interface-icon-pack/
    */
    div.folder {
      background-image: url("data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCA2NCA4MCIgeD0iMHB4IiB5PSIwcHgiPjxnPjxwYXRoIGQ9Ik01MC4xMywxOC4xN0gzMy44NGExLjcxLDEuNzEsMCwwLDEtMS40Ny0uODRMMzAsMTMuMTdhNC43MSw0LjcxLDAsMCwwLTQuMDYtMi4zNGgtMTJhNC43MSw0LjcxLDAsMCwwLTQuNyw0LjdWNDguNDdhNC43MSw0LjcxLDAsMCwwLDQuNyw0LjdINTAuMTNhNC43MSw0LjcxLDAsMCwwLDQuNy00LjdWMjIuODdBNC43MSw0LjcxLDAsMCwwLDUwLjEzLDE4LjE3Wm0xLjcsMzAuM2ExLjcsMS43LDAsMCwxLTEuNywxLjdIMTMuODdhMS43LDEuNywwLDAsMS0xLjctMS43VjE1LjUzYTEuNywxLjcsMCwwLDEsMS43LTEuN2gxMmExLjcxLDEuNzEsMCwwLDEsMS40Ny44NGwyLjQyLDQuMTVhNC43MSw0LjcxLDAsMCwwLDQuMDYsMi4zNEg1MC4xM2ExLjcsMS43LDAsMCwxLDEuNywxLjdaIi8+PC9nPjwvc3ZnPgo=");
    }
    div.file {
      background-image: url("data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCA2NCA4MCIgeD0iMHB4IiB5PSIwcHgiPjxnPjxwYXRoIGQ9Ik00OS44MiwyMC4xMywzOC42Niw4QTQuNTEsNC41MSwwLDAsMCwzNS4zNSw2LjVIMTcuNUE0LjUsNC41LDAsMCwwLDEzLDExVjUzYTQuNSw0LjUsMCwwLDAsNC41LDQuNWgyOUE0LjUsNC41LDAsMCwwLDUxLDUzVjIzLjE3QTQuNDksNC40OSwwLDAsMCw0OS44MiwyMC4xM1pNNDgsNTNhMS41LDEuNSwwLDAsMS0xLjUsMS41aC0yOUExLjUsMS41LDAsMCwxLDE2LDUzVjExYTEuNSwxLjUsMCwwLDEsMS41LTEuNUgzNS4zNWExLjUsMS41LDAsMCwxLDEuMTEuNDlMNDcuNjEsMjIuMTVhMS40OSwxLjQ5LDAsMCwxLC4zOSwxWiIvPjxwYXRoIGQ9Ik01MC43NiwyMmExLjUsMS41LDAsMCwxLTEuNSwxLjVoLTkuNmE0LjUsNC41LDAsMCwxLTQuNS00LjVWOC4zMWExLjUsMS41LDAsMCwxLDMsMFYxOWExLjUsMS41LDAsMCwwLDEuNSwxLjVoOS42QTEuNSwxLjUsMCwwLDEsNTAuNzYsMjJaIi8+PC9nPjwvc3ZnPgo=");
    }
    </style>
  </head>
  <body>
    <h1>{{.Path}}</h1>
    {{if (ne .Parent "")}}
      <div class="folder parent">
        <a href="{{.Parent}}">Parent Directory</a>
      </div>
    {{end}}
    {{range .Entries}}
      <div class="{{if .IsDir}}folder{{else}}file{{end}}">
        <a href="{{.Path}}">
          {{.Name}}{{if .IsDir}}/{{ end }}
        </a>
        <span class="size">{{ .Size }}</span>
      </div>
    {{end}}
  </body>
</html>
`))
}