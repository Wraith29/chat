package main

import (
	"chat/internal/consts"
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
		s, err := server.NewServer(consts.Host, consts.Port)

		if err != nil {
			fmt.Printf("Failed to create server: %+v\n", err)
			os.Exit(1)
		}
		defer s.Close()

		err = s.Listen()

		if err != nil {
			fmt.Printf("Failed to listen to server: %+v\n", err)
			os.Exit(1)
		}
	case "client":
		if len(os.Args) < 3 {
			fmt.Printf("Missing required argument 'name'\n")
			os.Exit(1)
		}

		name := os.Args[2]

		c, err := client.NewClient(consts.Host, consts.Port, name)

		if err != nil {
			fmt.Printf("Client creation failed: %+v\n", err)
			os.Exit(1)
		}
		defer c.Close()

		c.Run()
	default:
		fmt.Printf("invalid command '%s' expected one of ['server', 'client']\n", cmd)
		os.Exit(1)
	}
}
