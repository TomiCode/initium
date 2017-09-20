package app

import "log"
// import "strings"
import "net/http"

type Initium struct {
}

func init() {
  log.Println("Initium app global init.")
}

// Create application framework instance.
func Create() *Initium {
  return &Initium{}
}

// HTTP request handler.
func (app *Initium) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  log.Println("Request path:", r.URL.Path, "method:", r.Method)

  var handler = createHandler(w, r)
  log.Println(handler)

  var route = appRoutes.from(handler)
  if route == nil {
    log.Println("No route for", r.URL.Path)
    return
  }

  log.Println(route)
  log.Println(handler)
  route.methods[0].callback(handler)
}
