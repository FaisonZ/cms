package main

import (
	"fmt"
	"os"

	"faisonz.net/cms/web"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("You must provide a command")
		os.Exit(1)
	}

	switch args[0] {
	case "start":
		web.StartServer()
	case "db-init":
		fmt.Println("Needs to be implemented")
	default:
		fmt.Println("Invalid command:", args[0])
		os.Exit(1)
	}
}
