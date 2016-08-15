package main

import "fmt"

type BlogController struct {
}

func (controller* BlogController) routingRegister() []ControllerRoute {
  return []ControllerRoute{
    ControllerRoute{uri: "/", call: controller.index, template: "blog.index"},
    ControllerRoute{uri: "/add/{user}", call: controller.addPost},
  }
}

func (controller* BlogController) index(req *InitiumRequest) interface{} {
  fmt.Println("Index sees params:", req.vars)

  // if req.Middleware.User.IsLogged {
    // bleh.
  // }

  return []string{"abc", "bcd"}
}

func (controller* BlogController) addPost(req *InitiumRequest) interface{} {
  fmt.Println("BlogController addPost.")
  return nil
}
