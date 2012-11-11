package wiki

import (
  "os"
  "errors"
)

const pages = "pages/"

type GoWiki struct {
  Pages []Page
}

type Page struct {
  Title string
  Body string
}

func init() {
  fi, err := os.Stat(pages)
  if err != nil {
    panic(err)
  }

  if !fi.IsDir() {
    panic(errors.New("Wiki storage location is not a directory: " + pages))
  }
}

func Wiki() GoWiki {
  return GoWiki{Pages: []Page{Page{"page1", "Page 1"}, Page{"page2", "Page 2"}}}
}
