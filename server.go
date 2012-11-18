package main

import (
  "strings"
  "regexp"
  "log"
  "errors"
  "net/http"
  "html/template"
  "gowiki/wiki"
)

type GoWikiServer struct {
  *wiki.GoWiki
}

const views = "view/"
const mainView = "main.html"
const wikiView = "wiki.html"

var (
  templates = template.Must(template.ParseFiles(views + mainView, views + wikiView))
  validTitleRegExp = regexp.MustCompile("^/wiki/[^/.]+$")
  getTitleRegExp = regexp.MustCompile("[^/.]+$")
)

func (s *GoWikiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  var templateData interface{}
  var templateView string
  var err error

  log.Println("Serving request for " + r.URL.Path)

  switch {
  case strings.HasPrefix(r.URL.Path, "/wiki"):
    templateView = wikiView
    wikiTitle := getWikiTitle(r.URL.Path)
    log.Println("Request for wiki page: " + wikiTitle)

    if wikiTitle == "" {
      err = errors.New("Invalid Path to wiki")
    } else {
      templateData, err = s.GetWiki(wikiTitle)
    }
  case r.URL.Path == "/":
    templateView = mainView
    templateData, err = s.PageList()
  default:
    log.Println("Invalid Path: " + r.URL.Path)
    http.Error(w, "Not Found", http.StatusNotFound)
    return
  }

  if err == nil {
    err = templates.ExecuteTemplate(w, templateView, templateData)
  }

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}

func getWikiTitle(path string) (title string) {
  if validTitleRegExp.MatchString(path) {
    title = getTitleRegExp.FindString(path)
  }

  return
}


func main() {
  s := new(GoWikiServer)
  log.Println("Starting Server...")

  http.ListenAndServe(":8080", s)
}
