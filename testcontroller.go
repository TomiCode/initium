package main

import "log"

type TestController struct {
  App ApplicationInterface
}

func (controller *TestController) RegisterModule() *InitiumModule {
  return &InitiumModule{Title: "Test", RouteName: "test.index", ControllerAlias: "test_ctrl"}
}

func (controller *TestController) RegisterOptions() []*InitiumModuleCategory {
  return []*InitiumModuleCategory{
    &InitiumModuleCategory{
      Title: "TestCategory", 
      Options: []*ModuleOption{
        &ModuleOption{
          Name: "Test Option", 
          Route: "test.opt",
        },
      },
    },
  }
}

func (controller *TestController) RegisterRouting() []*ControllerRoute {
  return []*ControllerRoute{
    &ControllerRoute{uri: "/test/{doc}/{type}", call: controller.index, alias: "test.index", access: Permission_Auth_None},
  }
}

func (controller *TestController) index(request *InitiumRequest, params *RequestParameters) error {
  log.Println("[TestController] My parameters:", params)
  log.Println("[TestController] Value of 'doc':", params.GetValue("doc"))
  return nil
}