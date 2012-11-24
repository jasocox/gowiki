package wiki

type GoWiki struct { }

type WikiPage struct {
  Title string
  Body string
}

func (gw GoWiki) PageList() ([]WikiPage, error) {
  return getPageList()
}

func (gw GoWiki) GetWiki(title string) (WikiPage, error) {
  return getWiki(title)
}

func (gw GoWiki) CreateWiki(title string) (WikiPage, error) {
  return createOrUpdateWiki(title, "")
}

func (gw GoWiki) UpdateWiki(title string, body string) (WikiPage, error) {
  return createOrUpdateWiki(title, body)
}
