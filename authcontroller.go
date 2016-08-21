package main

import "fmt"

type AuthController struct {
  App ApplicationInterface
}

func (controller *AuthController) RoutingRegister() []ControllerRoute {
  return []ControllerRoute{
    ControllerRoute{uri: "/auth", call: controller.getLogin, name: "auth.login"},
    ControllerRoute{uri: "/auth", method: "POST", call: controller.postLogin, name: "auth.login.post"},
  }
}

func (controller *AuthController) getLogin(req *InitiumRequest) error {
  return controller.App.RenderTemplate(req, "auth.login", nil)
}

func (controller *AuthController) postLogin(req *InitiumRequest) error {
  fmt.Printf("postLogin: %+v\n", req)
  return nil
}
