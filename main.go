package main

import (
  "log"
  "net/http"
)

func main() {
  log.Println("Initium startup.")

  app := CreateInitium(true)
  app.LoadTemplates("templates")
  app.RegisterController(&BlogController{app})
  app.RegisterController(&AuthController{app})

  err := http.ListenAndServe("192.168.1.102:1234", app)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
