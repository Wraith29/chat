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
		log.Err("missing required param: 'server' or 'client'")
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch strings.ToLower(cmd) {
	case "server":
		s, err := server.NewServer(consts.Host, consts.Port)

		if err != nil {
			log.Err("server creation failed: %+v", err)
			os.Exit(1)
		}
		defer s.Close()

		err = s.Listen()

		if err != nil {
			log.Err("failed to listen to server: %+v", err)
			os.Exit(1)
		}
	case "client":
		if len(os.Args) < 3 {
			log.Err("missing required argument 'name'")
			os.Exit(1)
		}

		name := os.Args[2]

		c, err := client.NewClient(consts.Host, consts.Port, name)

		if err != nil {
			log.Err("client creation failed: %+v", err)
			os.Exit(1)
		}
		defer c.Close()

		go c.Run()
	default:
		log.Err("invalid command '%s' expected one of ['server', 'client']", cmd)
		os.Exit(1)
	}
}
