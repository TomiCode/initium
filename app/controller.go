package app

import "log"

// Every controller should inherit this struct.
type AppController struct {
  alias string
}

// Controller request method.
type RequestMethod func(bool) error

// Memory mapping for all controllers.
var appControllers *AppController[]

func init() {
  log.Println("Controller global init method.")
}

func (controller *AppController) Alias(alias string) (*AppController) {
  log.Println("Set controller alias to:", alias)
  controller.alias = alias
  return controller
}

func (controller *AppController) Register() bool {
  log.Printf("Register controller: %p\n", controller)
  return true
}

func (controller *AppController) hash() string {
  log.Println("Calculating hash for controller.")
  return "<none>"
}
