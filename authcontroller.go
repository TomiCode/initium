package main

import "log"

type AuthController struct {
  ApplicationInterface
  test int
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
    &ControllerRoute{uri: "/auth/form", method: "POST", call: controller.loginForm, alias: "auth.login.form"},
  }
}

func (controller *AuthController) getLogin(req *InitiumRequest, params *RequestParameters) error {
  return controller.RenderTemplate(req, "auth.login", nil)
}

func (controller *AuthController) postLogin(req *InitiumRequest, params *RequestParameters) error {
  var err = req.Request.ParseForm()
  if err != nil {
    log.Println("Can not parse form:", err)
    return nil
  }

  var user, pass = req.Request.Form.Get("email"), req.Request.Form.Get("passwd")

  err = controller.AuthenticateLogin(user, pass, req.Session)
  if err != nil {
    log.Println("Error occured while login:", err)
  }

  log.Println("Authenticate user:", user, pass)
  req.AddAlert("success", "Authorization", "Successful loggined in. Hello again!")
  return req.Redirect(controller.Route("blog.index"))
}

func (auth *AuthController) loginForm(req *InitiumRequest, params *RequestParameters) error {
  var response = InitiumHandler{}

  var err = req.Request.ParseForm()
  if err != nil {
    log.Println("Error while form parse:", err)
    return auth.RenderData(req, response)
  }
  log.Println(req.Request.Form);

  var user, pass = req.Request.Form.Get("email"), req.Request.Form.Get("passwd")
  err = auth.AuthenticateLogin(user, pass, req.Session)
  if err != nil {
    log.Println("Can not authorize user.")
    if auth.test % 2 == 0 {
      response.Error = "There was an error while authorizing your session. Seems, that You may have entered a wrong username or password."
    } else {
      response.Error = "Something unexpected happened."
    }
    auth.test++
    return auth.RenderData(req, response)
  }

  log.Println("Authorized user:", user)
  req.AddAlert("success", "Authorization", "Successful loggined in. Hello again!")

  response.Success = true
  response.Redirect = auth.Route("blog.index")
  return auth.RenderData(req, response)
}