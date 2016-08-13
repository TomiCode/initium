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

type RequestFunction func() bool

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

		expr, err := regexp.Compile(strings.Join(urlparts, "/"))
		if err != nil {
			log.Fatal("Can not compile regular expression for route:", v.path)
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
			log.Println("Found route for", r.URL.Path);
			route.fn()
		}
	}
}

