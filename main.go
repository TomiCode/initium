package main

// import (
//  "log"
//  "net/http"
// )

// func main() {
//   log.Println("Initium startup.")
// 
//   app := CreateInitium(true)
//   app.OpenDatabase("initium:123123@/initium_db")
//   defer app.CloseDatabase()
// 
//   app.LoadTemplates("templates")
// 
//   app.RegisterController(&AuthController{app, 0})
//   app.RegisterController(&BlogController{app})
//   app.RegisterController(&TestController{AppController{app}})
// 
//   err := http.ListenAndServe("localhost:1337", app)
//   if err != nil {
//     log.Fatal("ListenAndServe: ", err)
//   }
// }
