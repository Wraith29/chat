package main

import (
	"chat/internal/consts"
	"chat/pkg/client"
	"chat/pkg/server"

	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Missing required arg [\"server\" or \"client\"]\n")
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch strings.ToLower(cmd) {
	case "server":
		s, err := server.NewServer(consts.Host, consts.Port)

		if err != nil {
			fmt.Printf("Failed to create server: %+v\n", err)
			os.Exit(1)
		}
		defer s.Close()

		err = s.Listen()

		if err != nil {
			fmt.Printf("Failed to listen to server: %+v\n", err)
			os.Exit(1)
		}
	case "client":
		if len(os.Args) < 3 {
			fmt.Printf("Missing required argument 'name'\n")
			os.Exit(1)
		}

		name := os.Args[2]

		c, err := client.NewClient(consts.Host, consts.Port, name)

		if err != nil {
			fmt.Printf("Failed to create client: %+v\n", err)
			os.Exit(1)
		}

		app := tview.NewApplication()

		headerMsg := fmt.Sprintf("%s - Connected to %s:%s", name, consts.Host, consts.Port)

		header := tview.
			NewTextView().
			SetText(headerMsg).
			SetTextStyle(
				tcell.StyleDefault.
					Normal().
					Background(tcell.ColorGray).
					Foreground(tcell.ColorWhite))

		msgArea := tview.NewTextView().SetTextStyle(tcell.StyleDefault.Normal().Foreground(tcell.ColorWhite)).SetScrollable(true)

		go c.ReceiveInto(app, msgArea)

		msg := ""

		msgInputStyle := tcell.StyleDefault.Normal().Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)

		msgInput := tview.NewInputField()
		msgInput.
			SetPlaceholder("Enter a message...").
			SetChangedFunc(func(s string) { msg = s }).
			SetDoneFunc(func(k tcell.Key) {
				switch k {
				case tcell.KeyEnter:
					go c.Send(app, msgArea, msg)
					msg = ""
				}

				msgInput.SetText("")
			}).
			SetPlaceholderStyle(msgInputStyle).
			SetFieldStyle(msgInputStyle).
			SetLabelStyle(msgInputStyle)

		// Title (Client Name), Content (Messages), Input
		grid := tview.NewGrid().
			SetRows(1, 0, 1).
			AddItem(header, 0, 0, 1, 1, 0, 0, false).
			AddItem(msgArea, 1, 0, 1, 1, 0, 0, false).
			AddItem(msgInput, 2, 0, 1, 1, 0, 0, true)

		err = app.SetRoot(grid, true).Run()

		if err != nil {
			fmt.Printf("Failed to run application: %+v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("invalid command '%s' expected one of ['server', 'client']\n", cmd)
		os.Exit(1)
	}
}
