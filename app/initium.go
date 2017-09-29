package app

import "log"
// import "strings"
import "net/http"

type Initium struct {
  development bool
}

func init() {
  log.Println("Initium app global init.")
}

// Create application framework instance.
func Create(dev bool) *Initium {
  return &Initium{development: dev}
}

// HTTP request handler.
func (app *Initium) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  log.Println("Request path:", r.URL.Path, "method:", r.Method)

  var request = createRequest(r)
  log.Println("Local app request:", request)

  if app.development {
    if request.tryFile() {
      log.Println("Found file, serving content for this request.")
      return
    }
  }

  var route = appRoutes.from(request)
  if route == nil {
    log.Println("No route for", request.URL.Path)
    return
  }

  var callback = route.getCallback(request)
  log.Println("Callback:", callback)

  if callback != nil {
    response := callback(request)
    if response != nil {
      log.Println("Calling response callback..")
      response(nil)
    }
  }
}
