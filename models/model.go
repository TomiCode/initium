package models

import "log"

type AppModel interface {
  Create(interface{}) error
  Update(interface{}) error
  Destroy(interface{}) error
}


func init() {
  log.Println("Models init.")
}