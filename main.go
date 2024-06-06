package main

import (
	"chat/pkg/client"
	"chat/pkg/server"

	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Missing required arg [\"server\" or \"client\"]\n")
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch strings.ToLower(cmd) {
	case "server":
		app := server.CreateApp()
		defer app.Close()

		err := app.Run()

		if err != nil {
			panic(err)
		}
	case "client":
		if len(os.Args) < 3 {
			fmt.Printf("Missing required argument 'name'\n")
			os.Exit(1)
		}

		name := os.Args[2]

		app, err := client.CreateApp(name)

		if err != nil {
			fmt.Printf("Error creating app: %+v\n", err)
			os.Exit(1)
		}

		err = app.Run()

		if err != nil {
			fmt.Printf("Failed to run application: %+v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("invalid command '%s' expected one of ['server', 'client']\n", cmd)
		os.Exit(1)
	}
}
