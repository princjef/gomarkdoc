package main

import (
	"log"
)

func main() {
	log.SetFlags(0)

	cmd := buildCommand()

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
