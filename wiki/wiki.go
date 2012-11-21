package wiki

type GoWiki struct { }

type WikiPage struct {
  Title string
  Body string
}

func (gw GoWiki) PageList() (wp []WikiPage, err error) {
  return getPageList()
}

func (gw GoWiki) GetWiki(title string) (WikiPage, error) {
  return getWiki(title)
}

func (gw GoWiki) CreateWiki(title string, body string) (WikiPage, error) {
  return createWiki(title, body)
}
