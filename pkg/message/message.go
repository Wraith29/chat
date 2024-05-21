package message

type MessageType int

const (
	Connect = iota
	Send    = iota
)

type Message struct {
	Author  string `json:"author"`
	Message string `json:"message"`
}

func NewMessage(author, message string) Message {
	return Message{
		Author:  author,
		Message: message,
	}
}
