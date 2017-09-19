package app

import "log"
import "fmt"
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
  abstract string
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

// Create regular expression based on abstract routing path.
func (route *InternalRoute) regularPath() string {
  return fmt.Sprintf("^%s$", strings.Replace(route.abstract, "%v", "([^/]*?)", -1))
}

// Compile the routing path regular expression.
func (route *InternalRoute) compile() (err error) {
  route.path, err = regexp.Compile(route.regularPath())
  return
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
  if !strings.HasPrefix(route.path, "/") {
    log.Println("Creating local route prefix.")
    if controller, valid := appControllers[route.controller]; valid {
      route.path = fmt.Sprintf("/%s/%s", controller.alias, route.path)
    } else {
      log.Println("Invalid routing controller. Local routing requires a controller binding.")
      return false
    }
  }
  log.Println("Registering AppRoute", route.path, "into Initium.")

  // Extracting route and params into different variables.
  var params []string
  var parts []string = strings.Split(route.path, "/")

  for id, part := range(parts) {
    if strings.HasPrefix(part, ":") {
      log.Println("Found route parameter:", part)
      params = append(params, part[1:])
      parts[id] = "%v"
    }
  }

  var iroute = &InternalRoute{abstract: strings.Join(parts, "/")}
  if err := iroute.compile(); err != nil {
    log.Println("Error while routing compilation:", err)
    return false
  }

  return true
}
