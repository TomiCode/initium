package main

import "log"

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
  var err = req.Request.ParseForm()
  if err != nil {
    log.Println("Can not parse form:", err)
    return nil
  }

  var user, pass = req.Request.Form.Get("email"), req.Request.Form.Get("passwd")
  log.Println("Authenticate user:", user, pass)
  controller.App.AuthenticateUser(req, user, pass)

  return nil
}
