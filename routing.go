package main

import (
  "regexp"
  "strings"
  "log"
  "net/http"
  "reflect"
)

type RequestParameters struct {
}

type InitiumRequest struct {
  req *http.Request
  resw http.ResponseWriter

  params map[string]string
  template string
}

type RequestFunction func(*InitiumRequest)(bool)

type ControllerRoute struct {
  path string
  method string
  fn RequestFunction
}

type InitiumController interface {
  routingRegister() []ControllerRoute
}

type RoutingCollection struct {
  expr *regexp.Regexp
  params []string
  method string
  fn RequestFunction
}

type InitiumRouter struct {
  routes []*RoutingCollection
}

func (router* InitiumRouter) RegisterController(controller InitiumController) {
  for _, v := range controller.routingRegister() {
    urlparts := strings.Split(v.path, "/")

    var params []string
    for idx, part := range urlparts {
      if strings.HasPrefix(part, "{") {
        params = append(params, part[1:len(part) - 1])
        urlparts[idx] = "([^/]*?)"
      }
    }

    expr, err := regexp.Compile("^" + strings.Join(urlparts, "/") + "$")
    if err != nil {
      log.Println("[Warn] Can not compile regular expression for route:", v.path)
      continue;
    }
    router.routes = append(router.routes, &RoutingCollection{
      fn: v.fn,
      method: v.method,
      expr: expr,
      params: params,
    })

    log.Println("Registered", v.method, "route", v.path, "for", reflect.TypeOf(controller))
  }
}

func (router* InitiumRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  for _, route := range router.routes {
    if route.expr.MatchString(r.URL.Path) && route.method == r.Method {
      var params = make(map[string]string)

      if len(route.params) > 0 {
        // Check for [0] == r.URL.Path
        for idx, val := range route.expr.FindStringSubmatch(r.URL.Path) {
          if idx == 0 || val == "" {
            continue;
          }
          params[route.params[idx - 1]] = val
        }
      }
      for param, val := range r.URL.Query() {
        params[param] = val[0];
      }

      fnRequest := &InitiumRequest{req: r, resw: w, params: params}
      if !route.fn(fnRequest) {
        log.Println("Handling template result", fnRequest.template);
      }

      return
    }
  }
}

