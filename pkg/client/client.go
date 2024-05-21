package client

import (
	"chat/internal/consts"
	"chat/internal/log"
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

	_, err = conn.Write([]byte(name))

	if err != nil {
		return nil, err
	}

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
	go c.Write()
	go c.Receive()
}

func (c *client) Write() {
	stdin := os.Stdin

	buffer := make([]byte, consts.MessageSizeLimit)
	length, err := stdin.Read(buffer)

	if err != nil {
		if err != io.EOF {
			log.Err("Error: %+v", err)
			return
		}
	}

	msg := buffer[0:length]

	_, err = c.conn.Write(msg)

	if err != nil {
		log.Err("Error: %+v", err)
	}
}

func (c *client) readMessage() (string, error) {
	buffer := make([]byte, consts.MessageSizeLimit)
	length, err := c.conn.Read(buffer)

	if err != nil {
		return "", err
	}

	return string(buffer[0:length]), err
}

func (c *client) Receive() {
	msg, err := c.readMessage()
	log.Info("Received %s", msg)

	if err != nil {
		log.Err("Error: %+v", err)
		return
	}

	log.Info(msg)

}
