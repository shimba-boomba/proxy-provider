package main

import "encoding/json"
import "math/rand"
import "io/ioutil"
import "net/http"
import "strings"

const PROXY_LIST = "proxy-servers.txt"

type ProxyResponse struct {
  Proxy string
}

func handler(w http.ResponseWriter, r *http.Request) {
  bs, err := ioutil.ReadFile(PROXY_LIST)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  proxyServers := strings.Split(string(bs), "\n")
  randomProxy := proxyServers[rand.Intn(len(proxyServers))]

  resp := ProxyResponse{randomProxy}

  json, err := json.Marshal(resp)

  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(json)
}

func main() {
  http.HandleFunc("/random", handler)
  http.ListenAndServe(":3100", nil)
}
