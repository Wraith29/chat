package client

import (
	"bufio"
	"chat/internal/log"
	"chat/internal/message"
	"fmt"
	"io"
	"net"
	"os"
)

type client struct {
	StayOpen bool
	Name     string
	conn     net.Conn
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

	log.Info("connected")

	return &client{
		conn:     conn,
		StayOpen: true,
		Name:     name,
	}, nil
}

func (c *client) Close() {
	log.Info("Disconnecting from server")
	c.conn.Close()
}

func (c *client) Run() {
	for c.StayOpen {
		go c.Write()
		go c.Receive()
	}
}

func (c *client) Write() {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')

	if err != nil {
		log.Err("error reading input: %+v", err)
		return
	}

	msg := message.NewMessage(message.Send, c.Name, input)

	_, err = c.conn.Write(msg.ToBytes())

	if err != nil {
		log.Err("failed to write message: %+v", err)
		return
	}
}

func (c *client) readMessage() (message.Message, error) {
	reader := bufio.NewReader(c.conn)

	msg, err := reader.ReadString('\n')

	if err != nil {
		// Return a zeroed message, because we won't access it anyway
		return message.Message{}, err
	}

	return message.FromBytes([]byte(msg))
}

func (c *client) getAckMessage() message.Message {
	return message.NewMessage(message.Ack, c.Name, "")
}

func (c *client) Receive() {
	msg, err := c.readMessage()

	if err != nil && err != io.EOF {
		log.Err("failed to read message: %+v", err)
		return
	}

	log.Info("Received %s", msg.Message)

	ack := c.getAckMessage()

	_, err = c.conn.Write(ack.ToBytes())

	if err != nil {
		log.Err("failed to ack message: %+v", err)
		return
	}
}
