package wiki

import (
  "fmt"
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

func Wiki() (*GoWiki, error) {
  pageDir, err := os.Open(pages)
  if err != nil {
    return nil, err
  }

  wikiList, err := pageDir.Readdirnames(0)
  if err != nil {
    return nil, err
  }

  for page := range wikiList {
    fmt.Println("Wiki page file: " + pages + wikiList[page])
  }

  return &GoWiki{Pages: []Page{Page{"page1", "Page 1"}, Page{"page2", "Page 2"}}}, nil
}
