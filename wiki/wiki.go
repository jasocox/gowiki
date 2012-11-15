package wiki

import (
  "fmt"
  "regexp"
  "log"
  "os"
  "errors"
)

const pages = "pages/"
var validWikiFile = regexp.MustCompile("^[^/.]+[.]txt$")
var wikiTitleRegexp = regexp.MustCompile("[^.]+")

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
    log.Fatal(errors.New("Wiki storage location is not a directory: " + pages))
  }
}

func (gw GoWiki) PageList() (wp []WikiPage, err error) {
  pageDir, err := os.Open(pages)
  if err != nil {
    return
  }

  wikiList, err := pageDir.Readdirnames(0)
  if err != nil {
    return
  }

  for page := range wikiList {
    title, err := pageTitle(wikiList[page])
    if err != nil {
      log.Println(err)
      continue
    }
    fmt.Println("Wiki page title: " + title)
  }
  wp = []WikiPage{WikiPage{"page1", ""}, WikiPage{"page2", ""}}

  return wp, nil
}

func pageTitle(file string) (name string, err error) {
  if validWikiFile.MatchString(file) {
    name = wikiTitleRegexp.FindString(file)
    fmt.Println("name: " + name)
  } else {
    err = errors.New("Invalid wiki: " + file)
  }

  return
}
