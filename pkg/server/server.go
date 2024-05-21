package server

import (
	"chat/internal/log"
	"fmt"
	"net"
)

type server struct {
	listener    net.Listener
	connections []*Connection
	maxId       int
}

func NewServer(host, port string) (*server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))

	if err != nil {
		return nil, err
	}

	return &server{
		listener:    listener,
		connections: make([]*Connection, 0),
		maxId:       0,
	}, nil
}

func (s *server) Close() {
	s.listener.Close()
	for _, conn := range s.connections {
		conn.conn.Close()
	}
}

func (s *server) Listen() error {
	log.Info("Beginning Listen")

	for {
		conn, err := s.listener.Accept()

		if err != nil {
			log.Err("Error Accepting Connection: %+v\n", err)
			return err
		}

		connection, err := s.add(conn)

		if err != nil {
			return err
		}

		go connection.Handle(s)
	}
}

func (s *server) add(conn net.Conn) (*Connection, error) {
	result, err := NewConnection(s.maxId, conn)
	if err != nil {
		return nil, err
	}

	s.connections = append(s.connections, result)
	s.maxId++

	log.Info("Received connection from %s", result.name)

	return result, nil
}
