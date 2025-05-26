package twitch

import (
	"errors"
	"strings"
)

var ErrorNotValidFormat = errors.New("not valid format")

type Message struct {
	User    string
	Channel string
	Text    string
	Raw     string
}

// Format: ":username!user@user.tmi.twitch.tv PRIVMSG #channel :message text"
func ParseMessage(line string) (*Message, error) {
	parts := strings.SplitN(line, ":", 3)
	if len(parts) < 3 {
		return nil, ErrorNotValidFormat
	}

	meta := parts[1]
	text := parts[2]

	metaParts := strings.Fields(meta)
	if len(metaParts) < 3 {
		return nil, ErrorNotValidFormat
	}

	prefix := metaParts[0]
	channel := metaParts[2]
	user := strings.Split(prefix, "!")[0]
	return &Message{
		User:    strings.TrimSpace(user),
		Channel: strings.TrimSpace(channel),
		Text:    text,
		Raw:     line,
	}, nil
}
