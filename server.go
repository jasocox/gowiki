package main

import (
  "fmt"
  "net/http"
  "html/template"
)

type Server struct {}

// Wiki View
const views = "view/"
const mainView = "main.html"

var templates = template.Must(template.ParseFiles(views + mainView))

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  err := templates.ExecuteTemplate(w, mainView, nil)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

// Server
func main() {
  s := new(Server)
  fmt.Println("Starting Server...")

  http.ListenAndServe(":8080", s)
}
