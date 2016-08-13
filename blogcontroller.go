package main

import "fmt"

type BlogController struct {
}

func (controller* BlogController) routingRegister() []ControllerRoute {
  return []ControllerRoute{
    ControllerRoute{"/", "GET", controller.index},
    ControllerRoute{"/add/{user}", "GET", controller.addPost},
  }
}

func (controller* BlogController) index(req *InitiumRequest) bool {
  fmt.Println("Index sees params:", req.params)
  req.template = "blog_index"
  return false
}

func (controller* BlogController) addPost(req *InitiumRequest) bool {
  fmt.Println("BlogController addPost.")
  return true
}
