package server

import (
	"chat/internal/log"
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
	messages    chan message.Message
}

func NewServer(host, port string) (*server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))

	if err != nil {
		return nil, err
	}

	return &server{
		listener:    listener,
		connections: make(map[string]*Connection),
		mutex:       sync.Mutex{},
		messages:    make(chan message.Message),
		host:        host,
		port:        port,
	}, nil
}

func (s *server) Close() {
	s.listener.Close()
	for _, conn := range s.connections {
		conn.conn.Close()
	}
}

func (s *server) Listen() error {
	log.Info("Starting Server on %s:%s", s.host, s.port)

	// go s.broadcastMessages()

	for {
		conn, err := s.listener.Accept()

		if err != nil {
			log.Err("failed to accept connection: %+v\n", err)
			return err
		}

		s.mutex.Lock()
		connection, err := s.add(conn)
		s.mutex.Unlock()

		if err != nil {
			return err
		}

		go connection.Handle(s)
	}
}

// func (s *server) broadcastMessages() {
// 	log.Info("Broadcasting Messages")

// 	for {
// 		msg := <-s.messages
// 		log.Info("Message Pulled '%s'", msg.Message)

// 		s.mutex.Lock()
// 		for name, conn := range s.connections {
// 			if name == msg.Author {
// 				continue
// 			}

// 			log.Info("Writing %s to %s", msg.Message, name)

// 			_, err := conn.conn.Write(msg.ToBytes())

// 			if err != nil {
// 				log.Err("failed to write to connection: %+v", err)
// 			}
// 		}
// 		s.mutex.Unlock()
// 	}
// }

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

	log.Info("Received connection from %s", result.name)

	return result, nil
}
