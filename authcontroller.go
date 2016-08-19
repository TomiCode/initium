package main

// import "fmt"

type AuthController struct {
}

func (controller *AuthController) RoutingRegister() []ControllerRoute {
  return []ControllerRoute{
    ControllerRoute{uri: "/auth", call: controller.getLogin, template: "auth.login", name: "auth.login"},
  }
}

func (controller *AuthController) getLogin(req *InitiumRequest) interface{} {
  return true
}
