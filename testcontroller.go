package main

import "initium/app"
import "log"

type TestController struct {
  app.AppController
}

func init() {
  log.Println("TestController global init method.")

  controller := &TestController{}
  controller.Alias("test").Register()
  // controller.Alias("test")

  // app.CreateRoute("/", controller.index).Bind(controller).Register()
  // app.CerateRoute("test", controller.test).Bind(controller).Register()
}

func (c *TestController) index() {
  log.Println("TestController index method.")
}

func (c *TestController) test() {
  log.Println("TestController test method.")
}

// type AppController struct {
//  ApplicationInterface
//}

// type TestController struct {
//   AppController
// }

// func (controller *TestController) RegisterModule() *InitiumModule {
//  return &InitiumModule{Title: "Test", RouteName: "test.index", ControllerAlias: "test_ctrl"}
// }

// func (controller *TestController) RegisterOptions() []*InitiumModuleCategory {
//   return []*InitiumModuleCategory{
//     &InitiumModuleCategory{
//       Title: "TestCategory",
//       Options: []*ModuleOption{
//         &ModuleOption{
//           Name: "Test Option",
//           Route: "test.opt",
//         },
//         &ModuleOption{
//           Name: "Bleeeh",
//           Route: "test.index",
//         },
//       },
//     },
//   }
// }

// func (controller *TestController) RegisterRouting() []*ControllerRoute {
//   return []*ControllerRoute{
//     &ControllerRoute{uri: "/test/:doc", call: controller.index, alias: "test.index", access: Permission_NoAuth},
//     &ControllerRoute{uri: "/test/{abc}/{var}", call: controller.index, alias: "test.opt", access: Permission_Auth_None},
//   }
// }

// func (controller *TestController) index(request *InitiumRequest, params *RequestParameters) error {
//   log.Println("[TestController] My parameters:", params)
//   log.Println("[TestController] Value of 'doc':", params.GetValue("doc"))
//   return controller.RenderTemplate(request, "test.index", nil)
// }
