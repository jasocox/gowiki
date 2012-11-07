package wiki

type GoWiki struct {
  Pages []Page
}

type Page struct {
  Title string
  Body string
}

func Wiki(pages string) GoWiki {
  gw := GoWiki{Pages: []Page{Page{"page1", "Page 1"}, Page{"page2", "Page 2"}}}

  return gw
}
