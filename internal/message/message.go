package message

import (
	"errors"
)

type MessageType uint8

const (
	Connect = iota
	Send    = iota
	Ack     = iota
)

func (m MessageType) ToString() string {
	switch m {
	case Connect:
		return "Connect"
	case Send:
		return "Send"
	case Ack:
		return "Ack"
	default:
		return "Invalid"
	}
}

type Message struct {
	MessageType     MessageType
	Author, Message string
}

func NewMessage(messageType MessageType, author, message string) Message {
	return Message{
		MessageType: messageType,
		Author:      author,
		Message:     message,
	}
}

func NewAckMessage(author string) Message {
	return Message{
		MessageType: Ack,
		Author:      author,
		Message:     "",
	}
}

func (m *Message) ToBytes() []byte {
	buf := make([]byte, 0)

	buf = append(buf, byte(MessageType(m.MessageType)))

	authorLen := len(m.Author)
	buf = append(buf, []byte(string(rune(authorLen)))...)
	buf = append(buf, []byte(m.Author)...)

	msgLen := len(m.Message)
	buf = append(buf, []byte(string(rune(msgLen)))...)
	buf = append(buf, []byte(m.Message)...)

	buf = append(buf, '\n')

	return buf
}

type messageReader struct {
	bytes []byte
	index int
}

func newMessageReader(bytes []byte) messageReader {
	return messageReader{
		bytes: bytes,
		index: 0,
	}
}

func (m *messageReader) readByte() byte {
	b := m.bytes[m.index]

	m.index++

	return b
}

func (m *messageReader) readBytes(n int) []byte {
	b := m.bytes[m.index : m.index+n]

	m.index += n

	return b
}

func FromBytes(b []byte) (Message, error) {
	if len(b) < 7 {
		return Message{}, errors.New("invalid message, min length 7")
	}

	reader := newMessageReader(b)

	msgType := MessageType(reader.readByte())
	authorNameLen := int(reader.readByte())
	authorName := string(reader.readBytes(authorNameLen))
	msgLen := int(reader.readByte())
	msg := string(reader.readBytes(msgLen))

	return Message{
		MessageType: msgType,
		Author:      string(authorName),
		Message:     msg,
	}, nil
}
