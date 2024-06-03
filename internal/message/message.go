package message

import (
	"errors"
	"fmt"
	"time"
)

const timestampFmt = "02/01/06 15:04:05"

type MessageType uint8

const (
	Connect MessageType = iota
	Send                = iota
	Ack                 = iota
)

func (m MessageType) String() string {
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
	timestamp       time.Time
}

func NewMessage(messageType MessageType, author, message string) Message {
	return Message{
		MessageType: messageType,
		Author:      author,
		Message:     message,
		timestamp:   time.Now(),
	}
}

func (m *Message) time() string {
	return m.timestamp.Format(timestampFmt)
}

func (m *Message) String() string {
	return fmt.Sprintf("%s - %s: %s\n", m.time(), m.Author, m.Message)
}

func NewAckMessage(author string) Message {
	return NewMessage(Ack, author, "")
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

	msgTime := m.time()

	timeLen := len(msgTime)
	buf = append(buf, []byte(string(rune(timeLen)))...)
	buf = append(buf, []byte(msgTime)...)

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
	timeLen := int(reader.readByte())
	timestamp := string(reader.readBytes(timeLen))

	msgTime, err := time.Parse(timestampFmt, timestamp)

	if err != nil {
		return Message{}, err
	}

	return Message{
		MessageType: msgType,
		Author:      string(authorName),
		Message:     msg,
		timestamp:   msgTime,
	}, nil
}
