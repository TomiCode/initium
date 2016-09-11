package main

import "log"

type AuthController struct {
  App ApplicationInterface
}

func (controller* AuthController) RegisterModule() *InitiumModule {
  return &InitiumModule{Title: "Authorization"}
}

func (controller* AuthController) RegisterOptions() []*InitiumModuleCategory {
  return nil
}

func (controller* AuthController) RegisterRouting() []*ControllerRoute {
  return []*ControllerRoute{
    &ControllerRoute{uri: "/auth", call: controller.getLogin, alias: "auth.login", access: Permission_NoAuth},
    &ControllerRoute{uri: "/auth", method: "POST", call: controller.postLogin, alias: "auth.login.post"},
  }
}

func (controller *AuthController) getLogin(req *InitiumRequest, params *RequestParameters) error {
  return controller.App.RenderTemplate(req, "auth.login", nil)
}

func (controller *AuthController) postLogin(req *InitiumRequest, params *RequestParameters) error {
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
  req.AddAlert("success", "Authorization", "Successful loggined in. Hello again!")
  return req.Redirect(controller.App.Route("blog.index"))
}
