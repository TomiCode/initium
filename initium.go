package main

import (
  "log"
  "net/http"
)

func main() {
  log.Println("Starting..")
  router := &InitiumRouter{}

  router.RegisterController(&BlogController{})
  err := http.ListenAndServe(":1234", router)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
