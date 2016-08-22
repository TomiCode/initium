package main

import (
  "log"
  "database/sql"
  "math/rand"
  "fmt"
)

type InitiumUser struct {
  Name string
  Token int16
  Email string
  AuthToken string
}

func (request* InitiumRequest) IsAuthorized() bool {
  if request.User != nil {
    return true
  }
  return false
}

func (app* InitiumApp) StartAuthorization(request* InitiumRequest) error {
  log.Println("User authorization session started.")
  var database* sql.DB = app.GetDatabase()
  if database == nil {
    return CreateError("Database connection not valid", 901)
  }

  if !request.Session.IsValid("user_token") {
    return nil
  }
  var auth_token string = request.Session.GetValue("user_token").(string)
  var row = database.QueryRow("SELECT login, token, email FROM users WHERE auth_token=?", auth_token)
  log.Println("Authorization token:", auth_token)

  request.User = &InitiumUser{AuthToken: auth_token}
  var err = row.Scan(&request.User.Name, &request.User.Token, &request.User.Email)
  if err == sql.ErrNoRows {
    request.User = nil
    return nil
  }

  if err != nil {
    return err
  }
  return nil
}

func (app *InitiumApp) AuthenticateUser(req* InitiumRequest, user, pass string) bool {
  var database *sql.DB = app.GetDatabase()
  if database == nil {
    log.Println("Database connection problems.")
    return false
  }

  var row = database.QueryRow("SELECT id FROM users WHERE email=? AND password=?", user, pass)
  var user_id int

  var err = row.Scan(&user_id)
  if err == sql.ErrNoRows {
    log.Println("Username or password failed.")
    return false
  } else if err != nil {
    log.Println("Error while AuthenticateUser:", err)
    return false
  }

  var auth_id = make([]byte, 4)
  rand.Read(auth_id)

  var auth_string string = fmt.Sprintf("%x", auth_id)

  _, err = database.Exec("UPDATE users SET auth_token=? WHERE id=?", auth_string, user_id)
  if err != nil {
    log.Println("Error while updating user token:", err)
    return false
  }

  req.Session.SetValue("user_token", auth_string)
  return true
}

