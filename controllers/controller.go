package controllers

import "log"

const (
  MethodGET   = iota
  MethodPOST  = iota
)

type InitiumController interface{
  init() bool
}

type ReqFn func(bool, bool) error

type AccessFn func(bool) bool

type VisibleFn func(bool) bool

type InitiumRoute struct {
  uri string
  alias string
  method uint8
  callback ReqFn
  access AccessFn
}

type InitiumMenu struct {
  title string
  route int
  category int
  visible VisibleFn
}

func registerController(controller *InitiumController) {
  controller.init()
}

func registerRoute(route *InitiumRoute) int {
}

func registerMenuOption(menu *InitiumMenu) {
}

func registerMenuCategory(title string) int {
}

func init() {
  registerController(&BlogController{})
}