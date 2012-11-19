package wiki

import (
  "regexp"
  "log"
  "os"
  "io/ioutil"
  "errors"
)

const pages = "pages/"
var (
  validWikiFile = regexp.MustCompile("^[^/.]+[.]txt$")
  wikiTitleRegexp = regexp.MustCompile("[^.]+")
)

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
    wp = append(wp, WikiPage{title, ""})
  }

  return wp, nil
}

func (gw GoWiki) GetWiki(title string) (wp WikiPage, err error) {
  body, err := ioutil.ReadFile(pages + title + ".txt")
  wp.Title = title
  wp.Body = string(body)
  log.Println("Wiki body: " + wp.Body)

  return
}

func pageTitle(file string) (name string, err error) {
  if validWikiFile.MatchString(file) {
    name = wikiTitleRegexp.FindString(file)
  } else {
    err = errors.New("Invalid wiki: " + file)
  }

  return
}
