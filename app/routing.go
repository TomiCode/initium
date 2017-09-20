package app

import "log"
import "fmt"
import "regexp"
import "strings"
import "net/http"

// Application route element.
type AppRoute struct {
  path *regexp.Regexp
  alias string
  methods []RouteMethod
  abstract string
}

// Request method callback.
type RouteMethod struct {
  method MethodType
  callback RequestCallback
}

// Method type.
type MethodType uint8

// Route collection type.
type RouteCollection map[string]*AppRoute

// Request methods contants.
const (
  RequestGet    = 0x00
  RequestPost   = 0x01
  RequestPut    = 0x02
  RequestPatch  = 0x04
  RequestDelete = 0x08
)

// Application routing table.
var appRoutes RouteCollection

// Initialize the route mapping.
func init() {
  log.Println("Initializing route mapping.")
  appRoutes = make(RouteCollection)
}

// Get corresponding route for a request.
func (collection RouteCollection) get(request *http.Request) *AppRoute {
  for id, route := range collection {
    if route.path.MatchString(request.URL.Path) {
      log.Println("Found route for request:", id, route)
      return route
    }
  }
  return nil
}

// Create regular expression based on abstract routing path.
func (route *AppRoute) pathContent() string {
  return fmt.Sprintf("^%s$", strings.Replace(route.abstract, "%v", "([^/]*?)", -1))
}

// Compile the routing path regular expression.
func (route *AppRoute) compile() (err error) {
  route.path, err = regexp.Compile(route.pathContent())
  return
}

// Create routing object with parsed path and callback method.
func NewRoute(path string) *AppRoute {
  log.Println("Creating route path", path)

  // Extracting route and params into different variables.
  var params []interface{}
  var parts []string = strings.Split(path, "/")

  // Proceed through all the url path elements.
  for id, part := range(parts) {
    if strings.HasPrefix(part, ":") {
      log.Println("Found route parameter:", part)
      parts[id] = "%v"
      params = append(params, part[1:])
    }
  }

  // Return a new route instance.
  return &AppRoute{
    alias: fmt.Sprintf(strings.Join(parts[1:], "_"), params...),
    abstract: strings.Join(parts, "/"),
  }
}

// Normal basic get request for a route.
func (route *AppRoute) Get(callback RequestCallback) *AppRoute {
  route.methods = append(route.methods, RouteMethod{method: RequestGet, callback: callback})
  return route
}

// Post callback for route.
func (route *AppRoute) Post(callback RequestCallback) *AppRoute {
  route.methods = append(route.methods, RouteMethod{method: RequestPost, callback: callback})
  return route
}

// Put request callback.
func (route *AppRoute) Put(callback RequestCallback) *AppRoute {
  route.methods = append(route.methods, RouteMethod{method: RequestPut, callback: callback})
  return route
}

// Patch callback request method.
func (route *AppRoute) Patch(callback RequestCallback) *AppRoute {
  route.methods = append(route.methods, RouteMethod{method: RequestPatch, callback: callback})
  return route
}

// Delete callback request method.
func (route *AppRoute) Delete(callback RequestCallback) *AppRoute {
  route.methods = append(route.methods, RouteMethod{method: RequestDelete, callback: callback})
  return route
}

// Register the routing into Initium.
func (route *AppRoute) Register() (bool) {
  if route.alias == "" {
    route.alias = "root"
  }

  // Check if this route does not exists.
  if _, exists := appRoutes[route.alias]; exists {
    log.Println("Route", route.alias, "exists in application routing table!")
    return false
  }

  log.Println("Registering route", route.alias, "into Initium.")

  // Compile route regular expression.
  if err := route.compile(); err != nil {
    log.Println("Error while routing compilation:", err)
    return false
  }

  // Add compiled route to the application routing table.
  appRoutes[route.alias] = route
  return true
}
