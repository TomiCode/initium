package main

import (
  "log"
)

type InitiumUser struct {
  Name string
  Token int16
  Email string
}

type InitiumAuth struct {
  user* InitiumUser
}

type ApplicationAuth interface {
  CurrentUser() (*InitiumUser)
  IsAuthorized() bool
}

func (auth* InitiumAuth) StartAuthorization(app* ApplicationInterface, session* ApplicationSession) {
  log.Println("Hello world")
}
