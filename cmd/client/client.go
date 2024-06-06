package main

import (
	"chat/pkg/client"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Missing required argument \"name\"")

		os.Exit(1)
	}

	name := os.Args[1]

	app, err := client.CreateApp(name)

	if err != nil {
		fmt.Println("Error creating app: ", err)
		os.Exit(1)
	}
	defer app.Close()

	err = app.Run()

	if err != nil {
		fmt.Println("Error running app: ", err)
		os.Exit(1)
	}
}
