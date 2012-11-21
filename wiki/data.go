package wiki

import (
  "strings"
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

func init() {
  fi, err := os.Stat(pages)
  if err != nil {
    panic(err)
  }

  if !fi.IsDir() {
    log.Fatal(errors.New("Wiki storage location is not a directory: " + pages))
  }
}

func getPageList() (wps []WikiPage, err error) {
  wikiList, err := ioutil.ReadDir(pages)
  if err != nil {
    return
  }

  for page := range wikiList {
    title, err := pageTitle(wikiList[page].Name())
    if err != nil {
      log.Println(err)
      continue
    }
    wps = append(wps, WikiPage{title, ""})
  }

  return
}

func getWiki(title string) (wp WikiPage, err error) {
  body, err := ioutil.ReadFile(pages + title + ".txt")

  wp.Title = title
  wp.Body = strings.TrimSpace(string(body))

  return
}

func createWiki(title string, body string) (wp WikiPage, err error) {
  wp.Title = title
  wp.Body = strings.TrimSpace(body)

  err = ioutil.WriteFile(pages + title + ".txt", []byte(body), 0600)

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
