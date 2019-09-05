package main

import (
	"github.com/jeyldii/api/handlers"
	"log"
)

func main() {
	r := handlers.InitRouter()

	if err := r.Run(":5000"); err != nil {
		log.Fatal(err)
	}
}
