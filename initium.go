package main

import (
  "log"
  "net/http"
)

func main() {
  log.Println("Initium startup.")

  router := &InitiumRouter{}
  router.LoadTemplates("templates")
  router.RegisterController(&BlogController{})

  err := http.ListenAndServe(":1234", router)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
