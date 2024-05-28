package server

import (
	"bufio"
	"chat/internal/log"
	"chat/internal/message"
	"fmt"
	"io"
	"net"
	"strings"
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
		return nil, fmt.Errorf("invalid message type, expected \"Connect\" found %s", msg.MessageType.ToString())
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
	log.Info("Disconnected %s", c.name)
}

func (c *Connection) Handle(s *server) {
	defer c.Close(s)

	reader := bufio.NewReader(c.conn)
	msg, err := reader.ReadString('\n')

	if err != nil && err != io.EOF {
		log.Err("failed to read from conn: %+v", err)
		return
	}

	log.Info("raw msg %s", msg)

	msg = strings.Trim(msg, " \n\r")

	log.Info("Received %s from %s", msg, c.name)

	ack := message.NewAckMessage(c.name)

	_, err = c.conn.Write(ack.ToBytes())

	if err != nil {
		log.Err("failed to ack message: %+v", err)
		return
	}

	for name, conn := range s.connections {
		if c.name == name {
			continue
		}

		conn.conn.Write([]byte(msg))
	}
}
