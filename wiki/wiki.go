package wiki

type GoWiki struct { }

type WikiPage struct {
  Title string
  Body string
}

func (gw GoWiki) PageList() ([]WikiPage, error) {
  return getPageList()
}

func (gw GoWiki) GetWiki(title string) (wp WikiPage, err error) {
  wp, err = getWiki(title)
  if wp.Title == "" {
    wp.Title = title
  }

  return
}

func (gw GoWiki) CreateWiki(title string, body string) (WikiPage, error) {
  return createOrUpdateWiki(title, body)
}
