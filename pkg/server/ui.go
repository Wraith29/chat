package server

import (
	"chat/internal/consts"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type appState struct {
	app    *tview.Application
	server *server
	// Pointer to a string, so that it can be nil (No client selected)
	selectedClient *string
	clientList     *tview.List
	messageList    *tview.TextView
}

func CreateApp() appState {
	app := tview.NewApplication()

	messageList := tview.NewTextView()
	messageList.SetBorder(true).
		SetTitle("Messages").
		SetTitleAlign(tview.AlignLeft)

	clientList := tview.NewList().ShowSecondaryText(false)
	clientList.
		SetBorder(true).
		SetTitle("Connected Clients").
		SetTitleAlign(tview.AlignLeft)

	header := tview.NewTextView().
		SetText(fmt.Sprintf("Server running on %s:%s", consts.Host, consts.Port)).
		SetTextStyle(tcell.StyleDefault.
			Normal().
			Background(tcell.ColorGray).
			Foreground(tcell.ColorWhite))

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 1, 0, false).
		AddItem(tview.
			NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(clientList, 0, 2, true).
			AddItem(messageList, 0, 3, false), 0, 1, true)

	app.SetRoot(flex, true)

	server, err := NewServer(consts.Host, consts.Port, app, clientList, messageList)

	if err != nil {
		panic(err)
	}

	clientList.
		AddItem("All", "", 0, func() { server.refreshMessageList(true) }).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() != tcell.KeyRune {
				return event
			}

			char := event.Rune()

			if char == 'x' || char == 'X' {
				server.DisconnectCurrentUser()
			}

			return event
		})

	return appState{
		app:            app,
		server:         server,
		selectedClient: nil,
		clientList:     clientList,
		messageList:    messageList,
	}
}

func (a *appState) Close() {
	a.server.Close()
}

func (a *appState) Run() error {
	go a.server.Listen()

	return a.app.Run()
}
