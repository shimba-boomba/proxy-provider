package main

import "github.com/parnurzeal/gorequest"
import "io/ioutil"
import "strings"
import "time"
import "os"

const PROXY_LIST_SRC = "source-proxy-servers.txt"
const PROXY_LIST_DEST = "proxy-servers.txt"

const TEST_PAGE = "https://suggests.avia.yandex.ru/avia?query=Valencia&lang=ru"

const TIMEOUT_MS = 1000

func testServers(proxyServers []string, channel chan string) {
  for _, proxy := range proxyServers {
    testServer(proxy, channel)
  }
}

func testServer(proxy string, channel chan string) {
  request := gorequest.New()
  
  request.Proxy("http://" + proxy)
  request.Timeout(TIMEOUT_MS * time.Millisecond)

  request.Get(TEST_PAGE)

  _, _, err := request.End()

  if err != nil {
    channel <- ""
  } else {
    channel <- proxy
  }  
}

func main() {
  bs, err := ioutil.ReadFile(PROXY_LIST_SRC)

  if err != nil {
    return   
  }

  proxyServers := strings.Split(string(bs), "\n")
  proxyServersCnt := len(proxyServers)

  channel := make(chan string)

  start := 0
  chunk := 50

  parts := proxyServersCnt / chunk

  for i := 0; i <= parts; i++ {
    limit := start + chunk

    if i == parts {
      limit = proxyServersCnt
    }

    ch := proxyServers[start:limit]
    
    start += chunk

    go testServers(ch, channel)
  }

  file, err := os.Create(PROXY_LIST_DEST)
  
  if err != nil {
    return
  }
  
  defer file.Close()

  cnt := 0

  for {
    response := <- channel

    if len(response) > 0 {
      file.WriteString(response + "\n")
    }

    cnt += 1
    
    if cnt >= proxyServersCnt {
      break
    }
  }
}