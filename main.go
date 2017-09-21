package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"

	"github.com/euforic/kevast/kevast"
)

func main() {
	fmt.Println("Welcome to the kevast REPL")
	s := kevast.NewSession()
	log.Fatal(s.Run())
}
