package main

import "log"

type AuthController struct {
  App ApplicationInterface
}

func (controller* AuthController) PermissionNode() string {
  return "auth_ctrl"
}

func (controller* AuthController) RoutingRegister() []*ControllerRoute {
  return []*ControllerRoute{
    &ControllerRoute{uri: "/auth", call: controller.getLogin, alias: "auth.login"},
    &ControllerRoute{uri: "/auth", method: "POST", call: controller.postLogin, alias: "auth.login.post"},
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
  err = controller.App.AuthenticateLogin(user, pass, req.Session)
  if err != nil {
    log.Println("Error occured while login:", err)
  }
  log.Println("Authenticate user:", user, pass)
  return nil
}
