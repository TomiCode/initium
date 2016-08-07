package main

import (
  "fmt"
  "net/http"
)

type InitiumRouter struct {
  Test string;
}

func (p* InitiumRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  fmt.Printf("ServeHTTP: %+v\n", r);
}


