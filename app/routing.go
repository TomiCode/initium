package app

import "log"
import "fmt"
import "regexp"
import "strings"

// Local routing element.
type AppRoute struct {
  path string
  method uint8
  callback RequestCallback
  controller uint64
}

// Compiled routing element.
type InternalRoute struct {
  path *regexp.Regexp
  methods []RouteMethod
  abstract string
  controller uint64
}

// Request method callback.
type RouteMethod struct {
  method uint8
  callback RequestCallback
}

// Application routing table.
var appRoutes map[string]*InternalRoute

// Initialize the route mapping.
func init() {
  log.Println("Initializing route mapping.")
  appRoutes = make(map[string]*InternalRoute)
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

// Assign fields from AppRoute object.
func (route *InternalRoute) fromInstance(app_route *AppRoute) {
  if route.controller == 0 {
    route.controller = app_route.controller
  } else if route.controller != app_route.controller {
    log.Println("Routing controller mismatch! This shouldn't happen!")
  }

  route.methods = append(route.methods, RouteMethod{method: app_route.method, callback: app_route.callback})
  log.Println("Assigned from instance:", route.methods)
}

// Create routing object with parsed path and callback method.
func CreateRoute(path string, callback RequestCallback) *AppRoute {
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
    if controller, valid := appControllers[route.controller]; valid {
      route.path = fmt.Sprintf("/%s/%s", controller.alias, route.path)
    } else {
      log.Println("Invalid routing controller. Local routing requires a controller binding.")
      return false
    }
  }

  // Extracting route and params into different variables.
  var params []interface{}
  var parts []string = strings.Split(route.path, "/")

  for id, part := range(parts) {
    if strings.HasPrefix(part, ":") {
      log.Println("Found route parameter:", part)
      parts[id] = "%v"
      params = append(params, part[1:])
    }
  }


  // Create and compile the internal routing object.
  var alias string = fmt.Sprintf(strings.Join(parts, "_"), params...)
  var iroute = &InternalRoute{abstract: strings.Join(parts, "/")}
  if err := iroute.compile(); err != nil {
    log.Println("Error while routing compilation:", err)
    return false
  }

  log.Println("Compiled route params:", params, "as", alias)
  iroute.fromInstance(route)
  appRoutes[alias] = iroute

  return true
}
