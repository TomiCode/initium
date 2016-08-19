package main

import "fmt"

type InitiumSession struct {
}

type InitiumAuth struct {
  Session map[string]InitiumSession
}

func (auth *InitiumAuth) requestAuth() {
  fmt.Printf("\n")
}
