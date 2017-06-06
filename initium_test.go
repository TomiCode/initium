package main

import "testing"
import "log"

func TestGenerateHash(*testing.T) {
  initium := CreateInitium(true)
  if initium.GenerateHash("") != 0x1EEF {
    log.Fatal("Unknown hash generator.")
  }

  if initium.GenerateHash("auth.login.form") == initium.GenerateHash("auth.login.post") {
    log.Fatal("Hash failed.")
  }
  log.Println("Hash gen success.")
}