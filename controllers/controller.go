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

type BasicController struct {
  id int
}

func (base *BasicController) hash(alias string) {
  log.Println("Route alias to SetHash:", alias)
  base.id = 0
}

type InitiumController interface{
  init() bool
  hash(string)
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

  controller.hash(controller.alias())
  controller.init()

  selfController.used = false
  *controllers = append(*controllers, controller)
  return nil
}

type ReqFn func(bool, bool) error
type AccessFn func(bool) bool

type InitiumRoute struct {
  reg *regexp.Regexp
  method uint8
  params []string
  access AccessFn
  callback ReqFn
  controller int
}

type RoutingCollection []*InitiumRoute
type RoutingHelpers map[string]string

var controllers ControllerCollection
var routes RoutingCollection
var paths RoutingHelpers

func init() {
  log.Println("InitiumControllers package init.")
  paths = make(map[string]string)

  controllers.register(&BlogController{})
}

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

func registerRoute(controller int, uri string, callback ReqFn, method uint8, access AccessFn) int {
  var newroute *InitiumRoute = &InitiumRoute{controller: controller,
    callback: callback,
    method: method,
    access: access,
  }

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

  var alias string = routeAlias(parts, newroute.params)
  var path  string = strings.Join(parts, "/")
  if _, exist := paths[alias]; !exist {
    paths[alias] = path
  }

  newroute.reg, err = regexp.Compile("\\A/" + strings.Replace(path, "%v", "([\\w-]+?)", -1) + "/?\\z")
  if err != nil {
    log.Fatal("Error while regexp compilation:", err)
  }
  routes = append(routes, newroute)

  log.Println("Registered route:", alias, "id:", routeId)
  return routeId
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