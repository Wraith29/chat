package client

import (
	"bufio"
	"chat/internal/message"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type client struct {
	Name string
	conn net.Conn
	wg   sync.WaitGroup
}

func NewClient(host, port, name string) (*client, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))

	if err != nil {
		return nil, err
	}

	initMsg := message.NewMessage(message.Connect, name, name)

	_, err = conn.Write(initMsg.ToBytes())

	if err != nil {
		return nil, err
	}

	return &client{
		conn: conn,
		Name: name,
		wg:   sync.WaitGroup{},
	}, nil
}

func (c *client) Close() {
	fmt.Printf("Closing client %s\n", c.Name)
	c.conn.Close()
	c.wg.Wait()
}

func (c *client) Run() {
	c.wg.Add(2)

	go c.Write()
	go c.Receive()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	c.Close()
}

func (c *client) readFromServer() (message.Message, error) {
	reader := bufio.NewReader(c.conn)

	// rawMessage, err := reader.ReadBytes('\x00')
	rawMessage, err := reader.ReadBytes('\n')

	if err != nil {
		return message.Message{}, err
	}

	msg, err := message.FromBytes(rawMessage)

	if err != nil {
		return message.Message{}, err
	}

	return msg, nil
}

func (c *client) readFromStdin() (message.Message, error) {
	reader := bufio.NewReader(os.Stdin)

	// Need to write until a newline (enter) is given
	rawMessage, err := reader.ReadString('\n')

	if err != nil {
		return message.Message{}, err
	}

	msg := message.NewMessage(message.Send, c.Name, strings.Trim(rawMessage, " \n\r"))

	return msg, nil
}

func (c *client) Write() {
	defer c.wg.Done()

	for {
		msg, err := c.readFromStdin()

		if err != nil {
			return
		}

		_, err = c.conn.Write(msg.ToBytes())

		if err != nil {
			fmt.Printf("Failed to write to the server: %+v\n", err)
			return
		}
	}
}

func (c *client) Receive() {
	defer c.wg.Done()

	for {
		msg, err := c.readFromServer()

		if err != nil {
			fmt.Printf("Failed to read message from server: %+v\n", err)
			return
		}

		switch msg.MessageType {
		case message.Connect:
			return
		case message.Send:
			fmt.Printf("%s: %s\n", msg.Author, msg.Message)

			ack := message.NewAckMessage(c.Name)

			_, err = c.conn.Write(ack.ToBytes())

			if err != nil {
				fmt.Printf("Failed to acknowledge message: %+v\n", err)
				return
			}
		}
	}
}
