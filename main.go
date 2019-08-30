package main

import (
	"github.com/DeathHand/gateway/service"
	"log"
)

func main() {
	s, err := service.New(&service.Config{})
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(<-*s.ErrChan)
}
