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
  editView = "edit.html"
)

var (
  templates = template.Must(template.ParseFiles(views + mainView,
                                                views + wikiView,
                                                views + editView))
  validWikiUrl = regexp.MustCompile("^/wiki/[^/.]+$")
  validEditUrl = regexp.MustCompile("^/edit/[^/.]+$")
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
  case validWikiUrl.MatchString(r.URL.Path) && r.Method == "GET":
    templateView = wikiView
    wikiTitle := getWikiTitle(r.URL.Path)
    log.Println("Request for wiki page: " + wikiTitle)

    templateData, err = s.GetWiki(wikiTitle)

    if err != nil {
      log.Println("Page doesn't exist. Creating it?")
      http.Redirect(w, r, "/edit/" + wikiTitle, http.StatusFound)
    }
  case validWikiUrl.MatchString(r.URL.Path) && r.Method == "POST":
    templateView = wikiView
    wikiTitle := getWikiTitle(r.URL.Path)
    log.Println("Edited wiki page: " + wikiTitle)

    wikiBody := r.FormValue("body")
    templateData, err = s.CreateWiki(wikiTitle, wikiBody)

    if err == nil {
      r.Method = "GET"
      http.Redirect(w, r, "/wiki/" + wikiTitle, http.StatusFound)
      return
    }
  case validEditUrl.MatchString(r.URL.Path):
    templateView = editView
    wikiTitle := getWikiTitle(r.URL.Path)
    log.Println("Edit page: " + wikiTitle)

    templateData, _ = s.GetWiki(wikiTitle)
  case r.URL.Path == "/":
    templateView = mainView
    templateData, err = s.PageList()
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
