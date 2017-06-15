package controllers

import "log"
import "net/http"
import "regexp"
import "strings"

const (
  MethodGET     = iota
  MethodPOST    = iota
  MethodINVALID = 0xFF
)

type InitiumController interface{
  init() bool
  alias() string
}
type ControllerCollection []InitiumController

type ThisController struct {
  InitiumController
  used bool
}

var selfController ThisController

func (controllers *ControllerCollection) register(controller InitiumController) error {
  log.Println("Initializing controller", controller.alias())

  selfController.InitiumController = controller
  selfController.used = true
  controller.init()

  selfController.used = false
  *controllers = append(*controllers, controller)
  return nil
}

type ReqFn func(bool, bool) error
type AccessFn func(bool) bool
type VisibleFn func(bool) bool

type InitiumRoute struct {
  reg *regexp.Regexp
  method uint8
  params []string
  access AccessFn
  callback ReqFn
}

type RoutingCollection []*InitiumRoute
type RoutingHelpers map[string]string

var controllers ControllerCollection
var abstract_routes RoutingHelpers
var routes RoutingCollection

func routeAlias(parts []string, params []string) string {
  var param_id int = 0

  for index, part := range(parts) {
    if part == "%v" {
      parts[index] = params[param_id]
      param_id++
    }
  }
  return strings.Join(parts, "_")
}

func registerRoute(uri string, callback ReqFn, method uint8, access AccessFn) int {
  var newroute *InitiumRoute = &InitiumRoute{callback: callback, method: method, access: access}
  var routeId = len(routes)

  var parts []string = strings.Split(uri, "/")
  var err error

  if strings.Contains(uri, ":") {
    for index, part := range(parts) {
      if !strings.Contains(part, ":") {
        continue
      }
      parts[index] = "%v"
      newroute.params = append(newroute.params, part[1:])
    }
  }
  parts[0] = selfController.alias()

  var uriformat string = strings.Join(parts, "/")
  var alias string = routeAlias(parts, newroute.params)
  if _, exist := abstract_routes[alias]; !exist {
    abstract_routes[alias] = uriformat
  }

  newroute.reg, err = regexp.Compile("\\A/" + strings.Replace(uriformat, "%v", "([\\w-]+?)", -1) + "/?\\z")
  if err != nil {
    log.Fatal("Error while regexp compilation:", err)
  }
  routes = append(routes, newroute)

  log.Println("Registered route:", uri, "id:", routeId)
  return routeId
}

func registerMenuOption(title string, category int, visibile VisibleFn, route int) {
}

func registerMenuCategory(title string) int {
  return 0
}

func methodNumber(method string) uint8 {
  switch(method) {
  case "GET":
    return MethodGET
  case "POST":
    return MethodPOST
  default:
    return MethodINVALID
  }
}

func GetRoute(r *http.Request) (*InitiumRoute) {
  var method uint8 = methodNumber(r.Method)

  for _, route := range(routes) {
    if route.method != method {
      continue
    }
    if route.reg.MatchString(r.URL.Path) {
      return route
    }
  }
  return nil
}

func init() {
  log.Println("InitiumControllers package init.")
  abstract_routes = make(map[string]string)

  controllers.register(&BlogController{})
}