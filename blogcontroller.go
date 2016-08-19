package main

import "fmt"

type BlogController struct {
  App ApplicationInterface
}

func (controller* BlogController) RoutingRegister() []ControllerRoute {
  return []ControllerRoute{
    ControllerRoute{uri: "/", call: controller.index},
    ControllerRoute{uri: "/add/{user}", call: controller.addPost},
  }
}

func (controller* BlogController) index(req *InitiumRequest) (*InitiumError) {
  fmt.Println("Index sees params:", req.vars)

  // if req.Middleware.User.IsLogged {
    // bleh.
  // }

  controller.App.RenderTemplate(req, "blog.index", nil)
  return nil
}

func (controller* BlogController) addPost(req *InitiumRequest) (*InitiumError) {
  fmt.Println("BlogController addPost.")
  return nil
}
