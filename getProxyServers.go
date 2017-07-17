package main

import "github.com/parnurzeal/gorequest"
import "net/http"
import "regexp"
import "os"

const PROXY_HOST = "http://proxy-fresh.ru"

const SESSION_ID = "0c41f46cfa46d7a94bb7a0bc5499e9d7"
const PHPSESS_ID = "1deb659de8e5d42ee5006953e7bf7f10"

func main() {
  var _ string
  var errs []error
  var hash string
  var body string

  cookie := http.Cookie{ Name: "PHPSESSID", Value: PHPSESS_ID }

  hash, errs = getHash(&cookie)
  if errs != nil {
    panic(errs)
  }

  _, errs = checkSessId(&cookie, hash)
  if errs != nil {
    panic(errs)
  }

  _, errs = getLink(&cookie)
  if errs != nil {
    panic(errs)
  }

  body, errs = download(&cookie)
  if errs != nil {
    panic(errs)
  }

  file, err := os.Create("source-proxy-servers.txt")
  
  if err != nil {
    panic(err)
  }
  
  defer file.Close()
  file.WriteString(body)
}

func getHash(cookie *http.Cookie) (string, []error) {
  request := gorequest.New()

  path := "/proxy/type/https"
  request.Get(PROXY_HOST + path)

  request.AddCookie(cookie)

  _, body, err := request.End()

  if err != nil {
    return "", err
  }
  
  regExpr := regexp.MustCompile(`href="` + path + `/([0-9a-z]{32})/"`)
  hash := regExpr.FindStringSubmatch(body)[1]

  return hash, nil
}

func checkSessId(cookie *http.Cookie, hash string) (string, []error) {
  request := gorequest.New()

  path := "/ajax/index.php"
  request.Post(PROXY_HOST + path)

  request.Send("action=checksessid&sessid=" + SESSION_ID + "&hash=" + hash)

  request.AddCookie(cookie)

  _, body, err := request.End()

  if err != nil {
    return "", err
  }

  return body, nil
}

func getLink(cookie *http.Cookie) (string, []error) {
  request := gorequest.New()

  path := "/ajax/index.php"
  request.Post(PROXY_HOST + path)

  request.Send("action=getlink&sessid=" + SESSION_ID)

  request.AddCookie(cookie)

  _, body, err := request.End()

  if err != nil {
    return "", err
  }

  return body, nil
}

func download(cookie *http.Cookie) (string, []error) {
  request := gorequest.New()

  path := "/download"
  request.Post(PROXY_HOST + path)

  request.AddCookie(cookie)

  _, body, err := request.End()

  if err != nil {
    return "", err
  }

  return body, nil
}