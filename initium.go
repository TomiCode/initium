package main

import (
  "net/http"
  "log"
)

func main(){
  log.Println("Starting..")

  router := &InitiumRouter{}

  err := http.ListenAndServe(":1234", router)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
