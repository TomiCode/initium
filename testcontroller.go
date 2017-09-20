package main

import "initium/app"
import "log"

type TestController struct {
  app.AppController
}

func init() {
  log.Println("TestController global init method.")

  controller := &TestController{}
  controller.Register()

  app.NewRoute("/").Get(controller.index).Register()
  app.NewRoute("/test/:id").Get(controller.test).Register()
}

func (c *TestController) index(handler *app.Handler) error {
  log.Println("TestController index method.")
  return nil
}

func (c *TestController) test(handler *app.Handler) error {
  log.Println("TestController test method.")
  return nil
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
