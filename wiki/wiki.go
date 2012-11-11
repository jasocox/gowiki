package wiki

import "os"

const pages = "pages/"

type GoWiki struct {
  Pages []Page
}

type Page struct {
  Title string
  Body string
}

func init() {
  _, err := os.Stat(pages)
  if err != nil {
    panic(err)
  }
}

func Wiki() GoWiki {
  return GoWiki{Pages: []Page{Page{"page1", "Page 1"}, Page{"page2", "Page 2"}}}
}
