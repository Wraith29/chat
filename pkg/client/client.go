package client

import (
	"bufio"
	"chat/internal/message"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/rivo/tview"
)

type Client struct {
	Name string
	conn net.Conn
	open bool
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
		Name: name,
		conn: conn,
		open: true,
	}, nil
}

func (c *Client) Close(app *tview.Application) {
	c.open = false
	c.conn.Close()

	app.Stop()
	os.Exit(0)
}

func (c *Client) readFromServer() (message.Message, error) {
	reader := bufio.NewReader(c.conn)

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

func (c *Client) ReceiveInto(app *tview.Application, msgArea *tview.TextView) {
	for c.open {
		msg, err := c.readFromServer()

		if err != nil {
			panic(err)
		}

		switch msg.MessageType {
		case message.Connect:
			panic("Invalid message type")
		case message.Disconnect:
			go app.QueueUpdateDraw(func() {
				currentContent := msgArea.GetText(true)

				msgArea.SetText(currentContent + "Received a Disconnect Signal from the server. Closing connection")
				msgArea.ScrollToEnd()
			})
			time.Sleep(time.Second * 5)

			c.Close(app)
		case message.Send:
			go app.QueueUpdateDraw(func() {
				currentContent := msgArea.GetText(true)

				msgArea.SetText(currentContent + msg.String())
				msgArea.ScrollToEnd()
			})

			ackMsg := message.NewAckMessage(c.Name)

			_, err := c.conn.Write(ackMsg.ToBytes())

			if err != nil {
				fmt.Printf(":(")
				os.Exit(1)
			}
		}
	}
}

func (c *Client) executeCommand(app *tview.Application, command string) {
	if command == "quit" {
		c.Close(app)
	}
}
