package main

import (
  "fmt"
  "net/http"
  "html/template"
  "gowiki/wiki"
)

// Wiki Model
var gowiki = new(wiki.GoWiki)

// Wiki View
const views = "view/"
const mainView = "main.html"

var templates = template.Must(template.ParseFiles(views + mainView))

type Server struct {}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    return
  }

  pageList, err := gowiki.PageList()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  err = templates.ExecuteTemplate(w, mainView, pageList)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}

// Server
func main() {
  s := new(Server)
  fmt.Println("Starting Server...")

  http.ListenAndServe(":8080", s)
}
