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

func ServeWikiGet(w http.ResponseWriter, r *http.Request) {
  var (
    templateData interface{}
    err error
    status int
  )

  wikiTitle := mux.Vars(r)["title"]

  log.Println("ServeWikiGet: Request path: " + r.URL.Path + " Wiki page: " + wikiTitle)

  templateData, err = gowiki.GetWiki(wikiTitle)
  if err != nil {
    log.Println("Page doesn't exist. Creating it")
    http.Redirect(w, r, "/edit/" + wikiTitle, http.StatusFound)
  }

  rendorTemplateOrError(w, templateData, wikiView, err, status)
}

func ServeWikiPost(w http.ResponseWriter, r *http.Request) {
  var (
    templateData interface{}
    err error
    status int
  )

  wikiTitle := mux.Vars(r)["title"]

  log.Println("ServeWikiPost: Request path: " + r.URL.Path + " Wiki page: " + wikiTitle)

  wikiBody := r.FormValue("body")
  templateData, err = gowiki.UpdateWiki(wikiTitle, wikiBody)
  if err == nil {
    r.Method = "GET"
    http.Redirect(w, r, "/wiki/" + wikiTitle, http.StatusFound)
    return
  }

  rendorTemplateOrError(w, templateData, wikiView, err, status)
}

func ServeEdit(w http.ResponseWriter, r *http.Request) {
  var (
    templateData interface{}
    err error
    status int
  )

  wikiTitle := mux.Vars(r)["title"]

  log.Println("ServeEdit: Request path: " + r.URL.Path + " Wiki page: " + wikiTitle)

  templateData, err = gowiki.GetWiki(wikiTitle)
  if err != nil {
    templateData, err = gowiki.CreateWiki(wikiTitle)
  }

  rendorTemplateOrError(w, templateData, editView, err, status)
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
  var (
    templateData interface{}
    err error
    status int
  )

  log.Println("ServeHTTP: Request path: " + r.URL.Path)

  if r.URL.Path == "/" {
    templateData, err = gowiki.PageList()
  } else {
    log.Println("Invalid Path: " + r.URL.Path)
    err = errors.New("Not Found")
    status = http.StatusNotFound
  }

  rendorTemplateOrError(w, templateData, mainView, err, status)
}

func rendorTemplateOrError(w http.ResponseWriter,
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
  r.HandleFunc("/wiki/{title:[^/.]+}", ServeWikiGet).Methods("GET")
  r.HandleFunc("/wiki/{title:[^/.]+}", ServeWikiPost).Methods("POST")
  r.HandleFunc("/edit/{title:[^/.]+}", ServeEdit)
  r.HandleFunc("/", ServeHTTP)

  http.ListenAndServe(":8080", r)
}
