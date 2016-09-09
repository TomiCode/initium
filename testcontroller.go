package main

import "fmt"

type TestController struct {
  App ApplicationInterface
}

func (controller *TestController) RegisterModule() *InitiumModule {
  return &InitiumModule{Title: "Test", RouteName: "test.index", ControllerAlias: "test_ctrl"}
}

func (controller *TestController) RegisterOptions() []*InitiumModuleCategory {
  return []*InitiumModuleCategory{
    &InitiumModuleCategory{Title: "TestCategory", Options: []*ModuleOption{&ModuleOption{Name: "Test Option", Route: "test.opt"}}},
  }
}

func (controller *TestController) RegisterRouting() []*ControllerRoute {
  return []*ControllerRoute{
    &ControllerRoute{uri: "/test", call: controller.index, alias: "test.index"},
  }
}

func (controller *TestController) index(request *InitiumRequest) error {
  fmt.Println("Hello world")
  return nil
}