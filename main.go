package main

import (
	"chat/internal/consts"
	"chat/internal/log"
	"chat/pkg/client"
	"chat/pkg/server"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Err("Missing required param: 'server' or 'client'")
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch strings.ToLower(cmd) {
	case "server":
		s, err := server.NewServer(consts.Host, consts.Port)

		if err != nil {
			log.Err("Server Creation Failed: %+v", err)
			os.Exit(1)
		}
		defer s.Close()
		log.Info("Server running")

		s.Listen()
	case "client":
		if len(os.Args) < 3 {
			log.Err("Missing required argument 'name'")
			os.Exit(1)
		}

		name := os.Args[2]

		if len(name) > consts.ClientNameSizeLimit {
			log.Err("Client name is too long, max size: %d", consts.ClientNameSizeLimit)
			os.Exit(1)
		}

		c, err := client.NewClient(consts.Host, consts.Port, name)

		if err != nil {
			log.Err("Client Creation Failed: %+v", err)
			os.Exit(1)
		}
		defer c.Close()

		for c.StayOpen {
			go c.Run()

			if err != nil {
				c.StayOpen = false
				log.Err("Error - %+v", err)
			}
		}

	default:
		log.Err("Invalid command '%s' expected one of ['server', 'client']", cmd)
		os.Exit(1)
	}
}
