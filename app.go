package main

import "initium/app"
import "net/http"

func main() {
  log.Println("Web application startup.")

  var initium = app.Initium{}
  if err := http.ListenAndServe("127.0.0.1:1337", initium); err != nil {
    log.Fatal(err)
  }
}
