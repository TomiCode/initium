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

func (controller* BlogController) index(req *InitiumRequest) error {
  fmt.Println("Index sees params:", req.vars)

  // if req.Middleware.User.IsLogged {
    // bleh.
  // }
  return controller.App.RenderTemplate(req, "blog.index", nil)
}

func (controller* BlogController) addPost(req *InitiumRequest) error {
  fmt.Println("BlogController addPost.")
  return nil
}
