package main

import (
  "log"
  "net/http"
)

func main() {
  log.Println("Initium startup.")

  app := CreateInitium(true, "__initium_ssid", 16)
  app.OpenDatabase("initium:123123@/initium_db")
  defer app.CloseDatabase()

  app.LoadTemplates("templates")
  app.RegisterController(&BlogController{app})
  app.RegisterController(&AuthController{app})

  err := http.ListenAndServe("localhost:1337", app)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
