package main

import (
	"log"
	"net/http"
	"fmt"
)

type SuperController struct {
	Bleh string
}

func (p *SuperController) registerMethods() []TestFunc {
	aval := []TestFunc{p.testFunction}
	p.Bleh = "a123412"
	return aval
}

func (p *SuperController) testFunction(val int) bool {
	log.Println("Hello World from testFunction!", p.Bleh, val)
	p.Bleh = "123"
	return true
}

func main() {
	log.Println("Starting..")
  supcont := &SuperController{Bleh: "alamakota"}
	supcont.
	log.Printf("%+v, %+v\n", supcont, supcont.testFunction)
	router := &InitiumRouter{}

	router.RegisterController(supcont)
	log.Println(supcont.Bleh)
	err := http.ListenAndServe(":1234", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
