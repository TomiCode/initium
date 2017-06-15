package controllers

import "log"
import "http"
import "regexp"
import "strings"

const (
  MethodGET   = iota
  MethodPOST  = iota
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
  uri string
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
  var routeId = len(routes)

  if !strings.Contains(uri, ":") {
    routes = append(routes, &InitiumRoute{uri: uri, method: method, access: access, callback: callback})
    log.Println("Simple route registered:", uri, "id:", routeId)
    return routeId
  }

  var urlparts []string = strings.Split(uri, "/")
  var params []string

  for index, part := range urlparts {
    if !strings.Contains(part, ":") {
      continue
    }
    urlparts[index] = "%v"
    params = append(params, part[1:])
  }
  
  urlparts[0] = selfController.alias()

  var formatUri string = strings.Join(urlparts, "/")
  var alias string = routeAlias(urlparts, params)

  if _, exists := abstract_routes[alias]; !exists {
    abstract_routes[alias] = formatUri
  }

  regexp, err := regexp.Compile("^" + strings.Replace(formatUri, "%v", "([^/]*?)", -1) + "$")
  if err != nil {
    log.Fatal("Can not compile regular expression:", err)
  }

  routes = append(routes, &InitiumRoute{uri: uri, 
    reg: regexp,
    method: method,
    params: params,
    access: access,
    callback: callback,
  })

  log.Println("Registered expression route:", uri, "id:", routeId)
  return routeId
}

func registerMenuOption(title string, category int, visibile VisibleFn, route int) {
}

func registerMenuCategory(title string) int {
  return 0
}

func GetRouting(w http.ResponseWriter, r *http.Request) {
  
}

func init() {
  log.Println("InitiumControllers package init.")
  abstract_routes = make(map[string]string)

  controllers.register(&BlogController{})
}