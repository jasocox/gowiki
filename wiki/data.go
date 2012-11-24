package wiki

import (
  "strings"
  "regexp"
  "log"
  "errors"
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
)

const (
  dataURLs = "localhost"
  databaseID = "test"
  collectionID = "main"
)

var (
  session *mgo.Session
  database *mgo.Database
  mainWiki *mgo.Collection
  wikiTitleRegexp = regexp.MustCompile("[^.]+")
)

func init() {
  var err error

  err = connectDBSession()
  if err != nil {
    panic(err)
  }

  connectDBData()
}

func connectDBSession() (err error) {
  session, err = mgo.Dial(dataURLs)

  return
}

func connectDBData() {
  database = session.DB(databaseID)
  mainWiki = database.C(collectionID)
}

func checkConnectionStatusAndReconnect() {
  if session == nil {
    err := connectDBSession()
    if err == nil {
      connectDBData()
    }
  }
}

func getPageList() (wps []WikiPage, err error) {
  checkConnectionStatusAndReconnect()

  err = mainWiki.Find(nil).All(&wps)
  if err != nil {
    session = nil
    log.Println(err.Error())
    err = errors.New("Problems getting wiki list")
  }

  return
}

func getWiki(title string) (wp WikiPage, err error) {
  checkConnectionStatusAndReconnect()

  err = mainWiki.Find(bson.M{"title": title}).One(&wp)
  if err != nil {
    session = nil
    log.Println(err.Error())
    err = errors.New("Problems getting wiki named " + title)
  }

  return
}

func createOrUpdateWiki(title string, body string) (wp WikiPage, err error) {
  checkConnectionStatusAndReconnect()

  wp.Title = title
  wp.Body = strings.TrimSpace(body)

  selector := bson.M{"title": wp.Title}

  _, err = mainWiki.Upsert(selector, &wp)
  if err != nil {
    session = nil
    log.Println(err.Error())
    err = errors.New("Problems updating or creating wiki named " + title)
  }

  return
}
