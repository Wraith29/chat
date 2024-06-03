package server

import (
	"bufio"
	"chat/internal/message"
	"fmt"
	"io"
	"net"
)

type Connection struct {
	name string
	conn net.Conn
}

func NewConnection(conn net.Conn) (*Connection, error) {
	reader := bufio.NewReader(conn)
	rawMsg, err := reader.ReadBytes('\n')

	if err != nil && err != io.EOF {
		return nil, err
	}

	msg, err := message.FromBytes(rawMsg)

	if err != nil {
		return nil, err
	}

	if msg.MessageType != message.Connect {
		return nil, fmt.Errorf("invalid message type, expected \"Connect\" found %s", msg.MessageType.String())
	}

	return &Connection{
		name: msg.Author,
		conn: conn,
	}, nil
}

func (c *Connection) Close(s *server) {
	s.mutex.Lock()
	delete(s.connections, c.name)
	s.mutex.Unlock()
	c.conn.Close()
}

func (c *Connection) Handle(s *server) {
	defer c.Close(s)

	for {
		reader := bufio.NewReader(c.conn)

		rawMsg, err := reader.ReadBytes('\n')

		if err != nil && err != io.EOF {
			return
		}

		msg, err := message.FromBytes(rawMsg)

		if err != nil {
			return
		}

		switch msg.MessageType {
		case message.Connect:
			fmt.Printf("Unexpected message type \"Connect\" expected \"Send\" or \"Ack\"\n")
			return
		case message.Send:
			fmt.Printf("Received %s from %s\n", msg.Message, msg.Author)
			ack := message.NewAckMessage(c.name)
			_, err = c.conn.Write(ack.ToBytes())

			if err != nil {
				return
			}

			s.sendToAll(msg, c.name)
		case message.Ack:
			fmt.Printf("Received ack from %s\n", c.name)
		}

	}
}
