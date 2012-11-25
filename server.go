package main

import (
  "log"
  "errors"
  "net/http"
  "html/template"
  "gowiki/wiki"
  "github.com/gorilla/mux"
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
)

func ServeWiki(w http.ResponseWriter, r *http.Request) {
  var (
    templateData interface{}
    templateView string
    err error
    status int
  )

  wikiTitle := mux.Vars(r)["title"]
  templateView = wikiView

  log.Println("ServeWiki: Serving request for " + r.URL.Path)

  switch {
  case r.Method == "GET":
    log.Println("Request for wiki page: " + wikiTitle)

    templateData, err = gowiki.GetWiki(wikiTitle)

    if err != nil {
      log.Println("Page doesn't exist. Creating it?")
      http.Redirect(w, r, "/edit/" + wikiTitle, http.StatusFound)
    }
  case r.Method == "POST":
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

  handleErrorOrTemplate(w, templateData, templateView, err, status)
}

func ServeEdit(w http.ResponseWriter, r *http.Request) {
  var (
    templateData interface{}
    templateView string
    err error
    status int
  )

  wikiTitle := mux.Vars(r)["title"]
  templateView = editView

  log.Println("ServeEdit: Serving request for " + r.URL.Path)

  templateData, err = gowiki.GetWiki(wikiTitle)
  if err != nil {
    templateData, err = gowiki.CreateWiki(wikiTitle)
  }

  handleErrorOrTemplate(w, templateData, templateView, err, status)
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

  handleErrorOrTemplate(w, templateData, templateView, err, status)
}

func handleErrorOrTemplate(w http.ResponseWriter,
                           templateData interface{},
                           templateView string,
                           err error,
                           status int) {
  if (err != nil) && (status == 0) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  } else if err != nil {
    http.Error(w, err.Error(), status)
  } else {
    err = templates.ExecuteTemplate(w, templateView, templateData)
  }
}

func main() {
  gowiki = new(wiki.GoWiki)
  log.Println("Starting Server...")

  r := mux.NewRouter()
  r.HandleFunc("/wiki/{title:[^/.]+}", ServeWiki)
  r.HandleFunc("/edit/{title:[^/.]+}", ServeEdit)
  r.HandleFunc("/", ServeHTTP)

  http.ListenAndServe(":8080", r)
}
