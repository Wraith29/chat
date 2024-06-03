package client

import (
	"bufio"
	"chat/internal/message"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/rivo/tview"
)

type Client struct {
	Name string
	conn net.Conn
}

func NewClient(host, port, name string) (*Client, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))

	if err != nil {
		return nil, err
	}

	initMsg := message.NewMessage(message.Connect, name, name)

	_, err = conn.Write(initMsg.ToBytes())

	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
		Name: name,
	}, nil
}

func (c *Client) Close() {
	fmt.Printf("Closing client %s\n", c.Name)
	c.conn.Close()
}

func (c *Client) readFromServer() (message.Message, error) {
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

func (c *Client) readFromStdin() (message.Message, error) {
	reader := bufio.NewReader(os.Stdin)

	// Need to write until a newline (enter) is given
	rawMessage, err := reader.ReadString('\n')

	if err != nil {
		return message.Message{}, err
	}

	msg := message.NewMessage(message.Send, c.Name, strings.Trim(rawMessage, " \n\r"))

	return msg, nil
}

func (c *Client) Send(app *tview.Application, msgArea *tview.TextView, rawMsg string) {
	msg := message.NewMessage(message.Send, c.Name, rawMsg)

	_, err := c.conn.Write(msg.ToBytes())

	if err != nil {
		fmt.Printf("Failed to write to the server: %+v\n", err)
		return
	}

	app.QueueUpdateDraw(func() {
		currentContent := msgArea.GetText(true)
		msgArea.SetText(currentContent + msg.String())
	})
}

func (c *Client) Receive() {
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

func (c *Client) ReceiveInto(app *tview.Application, msgArea *tview.TextView) {
	for {
		msg, err := c.readFromServer()

		if err != nil {
			panic(err)
		}

		switch msg.MessageType {
		case message.Connect:
			panic("Invalid message type")
		case message.Send:
			go app.QueueUpdateDraw(func() {
				currentContent := msgArea.GetText(true)

				msgArea.SetText(currentContent + msg.String())
			})

			ackMsg := message.NewAckMessage(c.Name)

			_, err := c.conn.Write(ackMsg.ToBytes())

			if err != nil {
				panic(err)
			}
		}
	}
}
