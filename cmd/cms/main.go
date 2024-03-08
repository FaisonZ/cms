package main

import (
	"fmt"
	"os"

	"faisonz.net/cms/internal/db"
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
		if err := db.SetupDatabase(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		fmt.Println("Invalid command:", args[0])
		os.Exit(1)
	}
}
