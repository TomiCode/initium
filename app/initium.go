package app

import "log"
import "net/http"

type Initium struct {
}

func init() {
  log.Println("Initium app global init.")
}

func (app *Initium) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  log.Println("Call ServeHTTP method.")
}
