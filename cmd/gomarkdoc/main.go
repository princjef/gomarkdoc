package main

import (
	"fmt"
	"log"
	"os"
)

func init() {
	log.SetPrefix("gomarkdoc: ")
	log.SetFlags(0)
}

func main() {
	cmd := buildCommand()

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
