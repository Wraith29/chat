package server

import (
	"chat/internal/message"
	"errors"
	"fmt"
	"net"
	"sync"
)

type server struct {
	host, port  string
	listener    net.Listener
	connections map[string]*Connection
	mutex       sync.Mutex
}

func NewServer(host, port string) (*server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))

	if err != nil {
		return nil, err
	}

	fmt.Printf("Server started on %s:%s\n", host, port)

	return &server{
		listener:    listener,
		connections: make(map[string]*Connection),
		mutex:       sync.Mutex{},
		host:        host,
		port:        port,
	}, nil
}

func (s *server) Close() {
	fmt.Println("Closing Server")

	s.listener.Close()
	for _, conn := range s.connections {
		conn.Close(s)
	}

}

func (s *server) Listen() error {
	for {
		conn, err := s.listener.Accept()

		if err != nil {
			return err
		}

		s.mutex.Lock()
		connection, err := s.add(conn)
		s.mutex.Unlock()

		if err != nil {
			return err
		}

		fmt.Printf("Received connection from %s\n", connection.name)

		go connection.Handle(s)
	}
}

func (s *server) add(conn net.Conn) (*Connection, error) {
	result, err := NewConnection(conn)
	if err != nil {
		return nil, err
	}

	_, found := s.connections[result.name]

	if found {
		return nil, errors.New("Already have a connection with " + result.name)
	}

	s.connections[result.name] = result

	return result, nil
}

func (s *server) sendToAll(msg message.Message, author string) {
	s.mutex.Lock()

	for name, conn := range s.connections {
		if author == name {
			continue
		}

		_, err := conn.conn.Write(msg.ToBytes())

		if err != nil {
			return
		}
	}

	s.mutex.Unlock()
}
