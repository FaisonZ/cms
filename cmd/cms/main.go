package main

import (
	"fmt"
	"os"
	"strconv"

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
		if err := db.CreateMainDB(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "client-init":
		if len(args) != 2 {
			fmt.Println("Second argument must be a number greater than 0")
			os.Exit(1)
		}

		clientID, err := strconv.ParseInt(args[1], 10, 0)
		if err != nil || clientID < 1 {
			fmt.Println("Second argument must be a number greater than 0")
			os.Exit(1)
		}

		if err := db.CreateClientDB(int(clientID)); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		fmt.Println("Invalid command:", args[0])
		os.Exit(1)
	}
}
