package server

import (
	"chat/internal/message"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/rivo/tview"
)

type server struct {
	host, port     string
	listener       net.Listener
	connections    map[string]*Connection
	messages       []message.Message
	mutex          sync.Mutex
	app            *tview.Application
	clientList     *tview.List
	messageList    *tview.TextView
	selectedClient string
}

func NewServer(host, port string, app *tview.Application, clientList *tview.List, messageList *tview.TextView) (*server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))

	if err != nil {
		return nil, err
	}

	fmt.Printf("Server started on %s:%s\n", host, port)

	return &server{
		listener:       listener,
		connections:    make(map[string]*Connection),
		messages:       make([]message.Message, 0),
		mutex:          sync.Mutex{},
		host:           host,
		port:           port,
		app:            app,
		clientList:     clientList,
		messageList:    messageList,
		selectedClient: "",
	}, nil
}

func (s *server) Close() {
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

		go s.app.QueueUpdateDraw(func() {
			s.clientList.AddItem(connection.name, "", 0, func() {
				s.showMessagesFromUser(connection.name)
				s.selectedClient = connection.name
			})
		})

		go connection.Handle(s)
	}
}

func (s *server) DisconnectCurrentUser() {
	user, found := s.connections[s.selectedClient]

	if !found {
		return
	}

	dcMsg := message.NewMessage(message.Disconnect, "Server", "")

	_, err := user.conn.Write(dcMsg.ToBytes())

	if err != nil {
		panic(err)
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

func (s *server) refreshMessageList(resetTitle bool) {
	contents := strings.Builder{}
	for _, msg := range s.messages {
		contents.WriteString(msg.String())
	}

	go s.app.QueueUpdateDraw(func() {
		s.messageList.SetText(contents.String())

		if resetTitle {
			s.messageList.SetTitle("Messages")
		}
	})
}

func (s *server) showMessagesFromUser(clientName string) {
	messageDisplay := strings.Builder{}

	for _, msg := range s.messages {
		if msg.Author == clientName {
			messageDisplay.WriteString(msg.String())
		}
	}

	go s.app.QueueUpdateDraw(func() {
		s.messageList.SetText(messageDisplay.String())
		s.messageList.SetTitle("Messages - " + clientName)
	})
}
