package main

import (
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

const (
  views = "view/"
  mainView = "main.html"
  wikiView = "wiki.html"
)

var (
  templates = template.Must(template.ParseFiles(views + mainView, views + wikiView))
  validWikiUrl = regexp.MustCompile("^/wiki/[^/.]+$")
  getTitleRegExp = regexp.MustCompile("[^/.]+$")
)

func (s *GoWikiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  var (
    templateData interface{}
    templateView string
    err error
    status int
  )

  log.Println("Serving request for " + r.URL.Path)

  switch {
  case validWikiUrl.MatchString(r.URL.Path):
    templateView = wikiView
    wikiTitle := getWikiTitle(r.URL.Path)
    log.Println("Request for wiki page: " + wikiTitle)

    switch r.Method {
    case "GET":
      templateData, err = s.GetWiki(wikiTitle)
    default:
      err = errors.New("Method Not Allowed")
      status = http.StatusMethodNotAllowed
      log.Println("Attempt to " + r.Method + " to /wiki")
    }

    if err != nil && status != 0 {
      status = http.StatusInternalServerError
    }
  case r.URL.Path == "/":
    templateView = mainView
    templateData, err = s.PageList()

    if err != nil {
      log.Println("Failure to generate list of wiki pages")
      status = http.StatusInternalServerError
    }
  default:
    log.Println("Invalid Path: " + r.URL.Path)
    err = errors.New("Not Found")
    status = http.StatusNotFound
  }

  if (err != nil) && (status == 0) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  } else if err != nil {
    http.Error(w, err.Error(), status)
  } else {
    err = templates.ExecuteTemplate(w, templateView, templateData)
  }
}

func getWikiTitle(path string) (title string) {
  title = getTitleRegExp.FindString(path)

  return
}


func main() {
  s := new(GoWikiServer)
  log.Println("Starting Server...")

  http.ListenAndServe(":8080", s)
}
