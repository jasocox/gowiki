package main

import (
  "log"
  "net/http"
  "html/template"
  "gowiki/wiki"
)

type GoWikiServer struct {
  *wiki.GoWiki
}

const views = "view/"
const mainView = "main.html"

var templates = template.Must(template.ParseFiles(views + mainView))

func (s *GoWikiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    log.Println("Invalid Path: " + r.URL.Path)
    http.Error(w, "Not Found", http.StatusNotFound)
    return
  }

  log.Println("Serving request for " + r.URL.Path)

  pageList, err := s.PageList()
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

func main() {
  s := new(GoWikiServer)
  log.Println("Starting Server...")

  http.ListenAndServe(":8080", s)
}
