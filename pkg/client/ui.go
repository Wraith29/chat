package client

import (
	"chat/internal/consts"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type appState struct {
	client *Client
	app    *tview.Application
}

func CreateApp(name string) (*appState, error) {
	app := tview.NewApplication()

	c, err := NewClient(consts.Host, consts.Port, name)

	if err != nil {
		fmt.Printf("Error creating client: %+v\n", err)
		return nil, err
	}

	headerMsg := fmt.Sprintf("%s - Connected to %s:%s", name, consts.Host, consts.Port)

	header := tview.
		NewTextView().
		SetText(headerMsg).
		SetTextStyle(
			tcell.StyleDefault.
				Normal().
				Background(tcell.ColorGray).
				Foreground(tcell.ColorWhite))

	msgArea := tview.NewTextView().
		SetTextStyle(tcell.StyleDefault.
			Normal().
			Foreground(tcell.ColorWhite)).
		SetScrollable(true)

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
				if msg[0] == '/' {
					c.executeCommand(app, msg[1:])
					msg = ""
					break
				}

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

	app.SetRoot(grid, true)

	return &appState{
		client: c,
		app:    app,
	}, nil
}

func (a *appState) Run() error {
	return a.app.Run()
}

func (a *appState) Close() {
	a.client.Close(a.app)
}
