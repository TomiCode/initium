package main

import (
  "fmt"
  "net/http"
)

func init() {
  fmt.Println("Hello routing.go here!")
}

type RequestFunction func() bool

type ControllerRoute struct {
	path string
	fn   RequestFunction
}

type InitiumController interface {
}

type InitiumRouter struct {
	func routingRegister() []ControllerRoute
}

func (router* InitiumRouter) RegisterController(controller InitiumController) bool {
  intMethods := p.registerMethods();
  for _,v := range controller.routingRegister() {
    fmt.Printf("%+v\n", v);
    // v(123)
  }
  // fmt.Printf("%+v\n", intMethods);
  return false
}

func (p* InitiumRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  fmt.Printf("ServeHTTP: %+v\n", r)
}

