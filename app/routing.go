package app

import "log"
import "regexp"

// App routing element.
type AppRoute struct {
  expr *regexp.Regexp
  callback RequestMethod
  controller uint64
}

// Application routing table.
var appRoutes []*AppRoute

// Create routing object with parsed path and callback method.
func CreateRoute(path string, callback RequestMethod) (*AppRoute) {
  log.Println("Registering routing", path)
  return &AppRoute{}
}

// Change routing request method.
func (route *AppRoute) Method(method uint8) (*AppRoute) {
  log.Println("Change method to:", method)
  return route
}

// Bind routing to a application controller.
func (route *AppRoute) Bind(controller uint64) (*AppRoute) {
  log.Println("Binding route for controller:", controller)
  return route
}

// Register the routing into Initium.
func (route *AppRoute) Register() (bool) {
  log.Println("Registering route into Initium..")
  return true
}
