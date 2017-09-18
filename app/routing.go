package app

import "log"
import "regexp"

// App routing element.
type AppRoute struct {
  expr *regexp.Regexp
  callback RequestMethod
  controller uint64
}

func init() {
  log.Println("Routing global init method.")
}

func CreateRoute(path string, callback RequestMethod) (*AppRoute) {
  log.Println("Registering routing", path)
  return &AppRoute{}
}

func (route *AppRoute) Method(method uint8) (*AppRoute) {
  log.Println("Change method to:", method)
  return route
}

func (route *AppRoute) Bind(controller uint64) (*AppRoute) {
  log.Println("Binding route for controller:", controller)
  return route
}

func (route *AppRoute) Register() (bool) {
  log.Println("Registering route into Initium..")
  return true
}
