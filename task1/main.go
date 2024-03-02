package main

import (
	"log"
	"task1/server"
)

func main() {
	config := server.NewConfig()
	s := server.New(config)

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
