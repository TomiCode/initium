package main

// import "fmt"

type AuthController struct {
  App ApplicationInterface
}

func (controller *AuthController) RoutingRegister() []ControllerRoute {
  return []ControllerRoute{
    ControllerRoute{uri: "/auth", call: controller.getLogin, name: "auth.login"},
  }
}

func (controller *AuthController) getLogin(req *InitiumRequest) (*InitiumError) {

  return nil
}
