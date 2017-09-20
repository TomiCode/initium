package app

import "log"
import "net/http"

// Every controller should inherit this struct.
type AppController struct {
  alias string
  id uint64
}

// Application controller handler param.
type Handler struct {
  request struct {
    *http.Request
    http.ResponseWriter
  }
  raw_params []string
}

// Controller request method.
type RequestCallback func(*Handler) error

// Memory mapping for all controllers.
var appControllers map[uint64]*AppController

// Initialize appControllers mapping.
func init() {
  log.Println("Initializing controller mappings.")
  appControllers = make(map[uint64]*AppController)
}

// Create new internal request instance.
func createHandler(w http.ResponseWriter, r *http.Request) *Handler {
  return &Handler{request: struct{
      *http.Request
      http.ResponseWriter
    }{r, w},
  }
}

// Change controller namespace.
func (controller *AppController) Alias(alias string) (*AppController) {
  log.Println("Set controller alias to:", alias)
  controller.alias = alias
  return controller
}

// Register the controller into Intium.
func (controller *AppController) Register() bool {
  log.Printf("Register controller: %p\n", controller)
  appControllers[1] = controller
  controller.id = 1
  // appControllers = append(appControllers, controller)

  log.Println("Registered controllers:", len(appControllers))
  return true
}

// Controller id for binding and stuff.
func (controller *AppController) Id() uint64 {
  if controller.id == 0 {
    log.Println("Controller", controller.alias, "is not registered into Initium!")
  }
  return controller.id
}
