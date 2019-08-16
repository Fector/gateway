package main

import (
	"github.com/DeathHand/gateway/service"
	"log"
)

func main() {
	s, err := service.NewService()

	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(s.Run())
}
