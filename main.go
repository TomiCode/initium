package main

import (
  "log"
  "net/http"
)

func main() {
  log.Println("Initium startup.")

  application := &InitiumApp{}
  application.LoadTemplates("templates")
  application.RegisterController(&BlogController{})

  err := http.ListenAndServe(":1234", application)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
