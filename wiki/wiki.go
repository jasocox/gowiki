package wiki

import (
  "fmt"
  "os"
  "errors"
)

const pages = "pages/"

type GoWiki struct { }

type WikiPage struct {
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

func (gw GoWiki) PageList() ([]WikiPage, error) {
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

  return []WikiPage{WikiPage{"page1", ""}, WikiPage{"page2", ""}}, nil
}
