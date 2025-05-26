package twitch

import (
	"strings"
)

// MessageType enum
type MessageType int

const (
	MessageTypeUnknown MessageType = iota
	PRIVMSG
	USERNOTICE
	PING
	// add other types as needed
)

// Message interface
type Message interface {
	GetType() MessageType
	GetRaw() string
}

// Base struct
type MessageBase struct {
	Type MessageType
	Raw  string
}

// Implement Message interface
func (m MessageBase) GetType() MessageType {
	return m.Type
}

func (m MessageBase) GetRaw() string {
	return m.Raw
}

type MessagePing struct {
	MessageBase
}

type UnknowMessage struct {
	MessageBase
}

// :arzyy69!arzyy69@arzyy69.tmi.twitch.tv PRIVMSG #kapesnik69 :smazal twitch damn
type MessagePrivate struct {
	MessageBase
	Sender    string
	SenderTMI string
	Streamer  string
	Text      string
}

// :tmi.twitch.tv USERNOTICE #tartancz :adsger
type MessageNotice struct {
	MessageBase
	Streamer string
	Text     string
}

func parseMessage(line string) Message {
	line = strings.TrimSpace(line)

	if strings.HasPrefix(line, "PING") {
		return &MessagePing{
			MessageBase{
				Type: PING,
				Raw:  line,
			},
		}
	}
	
	parts := strings.Split(line, ":")

	if len(parts) < 3 {
		return &UnknowMessage{
			MessageBase{
				Type: MessageTypeUnknown,
				Raw:  line,
			},
		}
	}
	meta := strings.Split(parts[1], " ")
	if len(meta) < 3 {
		return &UnknowMessage{
			MessageBase{
				Type: MessageTypeUnknown,
				Raw:  line,
			},
		}
	}
	switch meta[1] {
	case "PRIVMSG":
		sender := strings.Split(parts[1], "!")
		if len(sender) < 2 {
			return &UnknowMessage{
				MessageBase{
					Type: MessageTypeUnknown,
					Raw:  line,
				},
			}
		}
		return &MessagePrivate{
			MessageBase: MessageBase{
				Type: PRIVMSG,
				Raw:  line,
			},
			Sender:    sender[0],
			SenderTMI: sender[1],
			Streamer:  meta[2],
			Text:      parts[2],
		}
	case "USERNOTICE":
		return &MessageNotice{
			MessageBase: MessageBase{
				Type: USERNOTICE,
				Raw:  line,
			},
			Streamer: meta[2],
			Text:     parts[2],
		}
	default:
		return &UnknowMessage{
			MessageBase{
				Type: MessageTypeUnknown,
				Raw:  line,
			},
		}
	}

}
