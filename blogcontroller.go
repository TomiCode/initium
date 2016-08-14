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

func (controller* BlogController) index(req *InitiumRequest) interface{} {
  fmt.Println("Index sees params:", req.params)
  req.template = "blog.index"
  return []string{"abc", "bcd"}
}

func (controller* BlogController) addPost(req *InitiumRequest) interface{} {
  fmt.Println("BlogController addPost.")
  return nil
}
