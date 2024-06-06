package main

import (
	"chat/pkg/server"
	"fmt"
	"os"
)

func main() {
	app := server.CreateApp()
	defer app.Close()

	err := app.Run()

	if err != nil {
		fmt.Println("Error running server: ", err)
		os.Exit(1)
	}
}
