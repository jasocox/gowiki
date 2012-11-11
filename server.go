package main

import (
  "fmt"
  "net/http"
  "html/template"
  "gowiki/wiki"
)

// Wiki Model

// Wiki View
const views = "view/"
const mainView = "main.html"

var templates = template.Must(template.ParseFiles(views + mainView))

type Server struct {}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  wikiPages, err := wiki.Wiki()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  err = templates.ExecuteTemplate(w, mainView, wikiPages)
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
