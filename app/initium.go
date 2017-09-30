package app

import "log"
// import "strings"
import "net/http"

type Initium struct {
  static string
}

func init() {
  log.Println("Initium app global init.")
}

// Create application framework instance.
func Create() *Initium {
  return &Initium{}
}

func (app *Initium) EnableStaticFiles(dir string) *Initium {
  log.Println("Enabling static file serving from", dir)
  app.static = dir
  return app
}

// HTTP request handler.
func (app *Initium) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  log.Println("Request path:", r.URL.Path, "method:", r.Method)

  var request = createRequest(r)
  log.Println("Local app request:", request)

  if app.static != "" {
    log.Println("Static files enabled, checking request for a static file..")
    if app.tryServeFile(request, w) {
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
      var handler = createHandler(w)
      response(handler)
    }
  }
}

func (app *Initium) tryServeFile(request *Request, w http.ResponseWriter) bool {
  return false
}
