package app

import "log"

// Every controller should inherit this struct.
type AppController struct {
  alias string
}

// Controller request method.
type RequestMethod func(bool) error

// Memory mapping for all controllers.
var appControllers []*AppController

// Change controller namespace.
func (controller *AppController) Alias(alias string) (*AppController) {
  log.Println("Set controller alias to:", alias)
  controller.alias = alias
  return controller
}

// Register the controller into Intium.
func (controller *AppController) Register() bool {
  log.Printf("Register controller: %p\n", controller)
  appControllers = append(appControllers, controller)

  log.Println("Registered controllers:", len(appControllers))
  return true
}
