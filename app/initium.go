package app

import "log"
import "strings"
import "net/http"

type Initium struct {
}

func init() {
  log.Println("Initium app global init.")
}

func (app *Initium) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  log.Println("Request path:", r.URL.Path, "method:", r.Method)
  if r.Method == "GET" && strings.Contains(r.URL.Path, "assets") {
    log.Println("Handling asset file request.")
    http.ServeFile(w, r, "public" + r.URL.Path)
    return
  }

  route := appRoutes.get(r)
  log.Println(route)
}
