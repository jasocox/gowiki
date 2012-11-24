package main

import (
  "regexp"
  "log"
  "errors"
  "net/http"
  "html/template"
  "gowiki/wiki"
)

const (
  views = "view/"
  mainView = "main.html"
  wikiView = "wiki.html"
  editView = "edit.html"
)

var (
  gowiki *wiki.GoWiki
  templates = template.Must(template.ParseFiles(views + mainView,
                                                views + wikiView,
                                                views + editView))
  validWikiUrl = regexp.MustCompile("^/wiki/[^/.]+$")
  validEditUrl = regexp.MustCompile("^/edit/[^/.]+$")
  getTitleRegExp = regexp.MustCompile("[^/.]+$")
)

func ServeWiki(w http.ResponseWriter, r *http.Request) {
  var (
    templateData interface{}
    templateView string
    err error
    status int
  )

  log.Println("ServeWiki: Serving request for " + r.URL.Path)

  switch {
  case validWikiUrl.MatchString(r.URL.Path) && r.Method == "GET":
    templateView = wikiView
    wikiTitle := getWikiTitle(r.URL.Path)
    log.Println("Request for wiki page: " + wikiTitle)

    templateData, err = gowiki.GetWiki(wikiTitle)

    if err != nil {
      log.Println("Page doesn't exist. Creating it?")
      http.Redirect(w, r, "/edit/" + wikiTitle, http.StatusFound)
    }
  case validWikiUrl.MatchString(r.URL.Path) && r.Method == "POST":
    templateView = wikiView
    wikiTitle := getWikiTitle(r.URL.Path)
    log.Println("Edited wiki page: " + wikiTitle)

    wikiBody := r.FormValue("body")
    templateData, err = gowiki.UpdateWiki(wikiTitle, wikiBody)

    if err == nil {
      r.Method = "GET"
      http.Redirect(w, r, "/wiki/" + wikiTitle, http.StatusFound)
      return
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

func ServeEdit(w http.ResponseWriter, r *http.Request) {
  var (
    templateData interface{}
    templateView string
    err error
    status int
  )

  log.Println("ServeEdit: Serving request for " + r.URL.Path)

  templateView = editView
  wikiTitle := getWikiTitle(r.URL.Path)
  log.Println("Edit page: " + wikiTitle)

  templateData, err = gowiki.GetWiki(wikiTitle)
  if err != nil {
    templateData, err = gowiki.CreateWiki(wikiTitle)
  }

  if (err != nil) && (status == 0) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  } else if err != nil {
    http.Error(w, err.Error(), status)
  } else {
    err = templates.ExecuteTemplate(w, templateView, templateData)
  }
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
  var (
    templateData interface{}
    templateView string
    err error
    status int
  )

  log.Println("ServeHTTP: Serving request for " + r.URL.Path)

  if r.URL.Path == "/" {
    templateView = mainView
    templateData, err = gowiki.PageList()
  } else {
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
  gowiki = new(wiki.GoWiki)
  log.Println("Starting Server...")

  http.HandleFunc("/wiki/", ServeWiki)
  http.HandleFunc("/edit/", ServeEdit)
  http.HandleFunc("/", ServeHTTP)
  http.ListenAndServe(":8080", nil)
}
