package server

import (
	"chat/internal/consts"
	"chat/internal/log"
	"chat/internal/utils"
	"net"
)

type Connection struct {
	id   int
	name string
	conn net.Conn
}

func NewConnection(id int, conn net.Conn) (*Connection, error) {
	buffer := make([]byte, consts.ClientNameSizeLimit)
	length, err := conn.Read(buffer)

	if err != nil {
		return nil, err
	}

	name := buffer[0:length]

	return &Connection{
		id:   id,
		name: string(name),
		conn: conn,
	}, nil
}

func (c *Connection) Handle(s *server) {
	buffer := make([]byte, consts.MessageSizeLimit)
	length, err := c.conn.Read(buffer)

	if err != nil {
		log.Err("Error: %+v", err)
		return
	}

	msg := utils.SanitiseMessage(buffer[0:length])

	log.Info("Received %s from %s", string(msg), c.name)

	for _, conn := range s.connections {
		if conn.id == c.id {
			continue
		}

		log.Info("Writing %s to %s", string(msg), conn.name)
		_, err = conn.conn.Write(msg)

		if err != nil {
			log.Err("Error: %+v", err)
		}
	}
}
