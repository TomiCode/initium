package app

import "log"
import "regexp"
import "strings"

// Local routing element.
type AppRoute struct {
  path string
  method uint8
  callback RequestMethod
  controller uint64
}

// Compiled routing element.
type InternalRoute struct {
  path *regexp.Regexp
  methods []uint8
  callback RequestMethod
  controller uint64
}

// Application routing table.
// var appRoutes []*AppRoute
var appRoutes map[uint64]*InternalRoute

// Initialize the route mapping.
func init() {
  log.Println("Initializing route mapping.")
  appRoutes = make(map[uint64]*InternalRoute)
}

// Create routing object with parsed path and callback method.
func CreateRoute(path string, callback RequestMethod) *AppRoute {
  log.Println("Creating route path", path)
  return &AppRoute{path: path, callback: callback}
}

// Change routing request method.
func (route *AppRoute) Method(method uint8) (*AppRoute) {
  log.Println("Route", route.path, "changed to method:", method)
  route.method = method
  return route
}

// Bind routing to a application controller.
func (route *AppRoute) Bind(controller uint64) (*AppRoute) {
  log.Println("Bind route", route.path, "for controller:", controller)
  route.controller = controller
  return route
}

// Register the routing into Initium.
func (route *AppRoute) Register() (bool) {
  log.Println("Registering AppRoute", route.path, "into Initium.")
  if !strings.HasPrefix(route.path, "/") {
    log.Println("Creating local route prefix.")

  }
  return true
}
