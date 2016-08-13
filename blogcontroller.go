package main

import "fmt"

type BlogController struct {
}

func (controller* BlogController) routingRegister() []ControllerRoute {
  return []ControllerRoute{
    ControllerRoute{"/{index}", "GET", controller.index},
  }
}

func (controller* BlogController) index() bool {
  fmt.Println("BlogController index.")
  return true
}
