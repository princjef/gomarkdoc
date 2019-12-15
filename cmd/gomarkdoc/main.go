package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	cmd := buildCommand()

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
