package twitch

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

const (
	ircServer = "irc.chat.twitch.tv:6667"
)

var (
	ErrTooMuchStreamers = errors.New("Client can handle only 50 streamers")
)

type Client struct {
	oauth         string
	nick          string
	streamers     []string
	authenticated bool

	conn net.Conn

	reader *bufio.Reader

	onChatMessage func(m *MessagePrivate)

	onChatNotice func(m *MessageNotice)

	onPingPong func(m *MessagePing)

	onAnyMessage func(m Message)

	onUnknowMessage func(m *UnknowMessage)
}

func NewClient(oauth, nick string, streamers ...string) *Client {
	return &Client{
		oauth:     oauth,
		nick:      nick,
		streamers: streamers,
	}
}

func NewAnonymousClient(streamers ...string) *Client {
	return NewClient("oauth:59301", "justinfan123123", streamers...)
}

func (c *Client) SetOnChatMessage(callback func(m *MessagePrivate)) {
	c.onChatMessage = callback
}

func (c *Client) SetOnPingPong(callback func(m *MessagePing)) {
	c.onPingPong = callback
}

func (c *Client) SetOnChatNotice(callback func(m *MessageNotice)) {
	c.onChatNotice = callback
}

func (c *Client) SetOnAnyMessage(callback func(m Message)) {
	c.onAnyMessage = callback
}

func (c *Client) SetOnUnknowMessage(callback func(m *UnknowMessage)) {
	c.onUnknowMessage = callback
}

func (c *Client) Listen() error {
	for {
		if err := c.connectAndJoin(); err != nil {
			return fmt.Errorf("failed while connecting: %w", err)
		}
		for {
			line, err := c.reader.ReadString('\n')
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return fmt.Errorf("failed to read from connection: %w", err)
			}
			c.handleLine(line)
		}
	}
}

func (c *Client) handleLine(line string) {
	message := parseMessage(line)
	if c.onAnyMessage != nil {
		go c.onAnyMessage(message)
	}

	switch msg := message.(type) {
	case *UnknowMessage:
		if c.onUnknowMessage != nil {
			go c.onUnknowMessage(msg)
		}
	case *MessagePing:
		if c.onPingPong != nil {
			go c.onPingPong(msg)
		}
		c.SendPong(msg.GetRaw())
	case *MessageNotice:
		if c.onChatNotice != nil {
			go c.onChatNotice(msg)
		}
	case *MessagePrivate:
		if c.onChatMessage != nil {
			go c.onChatMessage(msg)
		}
	}
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) connectAndJoin() error {
	if len(c.streamers) == 0 {
		return errors.New("nothing to listen (c.streamers is empty)")
	}
	if err := c.makeConnection(); err != nil {
		return fmt.Errorf("failed to connect to Twitch IRC: %w", err)
	}
	if err := c.authenticate(); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	c.makeJoins()
	return nil
}

func (c *Client) makeConnection() error {
	conn, err := net.Dial("tcp", ircServer)
	if err != nil {
		return fmt.Errorf("failed to connect to Twitch IRC: %w", err)
	}

	c.conn = conn
	c.reader = bufio.NewReader(conn)

	return nil
}

func (c *Client) authenticate() error {
	c.conn.SetDeadline(time.Now().Add(30 * time.Second))
	defer c.conn.SetDeadline(time.Time{})

	fmt.Fprintf(c, "PASS %s\r\n", c.oauth)
	fmt.Fprintf(c, "NICK %s\r\n", c.nick)

	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading from server: %v", err)
		}
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Login authentication failed") {
			return errors.New("login authentication failed")
		} else if strings.Contains(line, "001 "+c.nick) {
			log.Println("✅ Authentication successful! Connected to Twitch IRC.")
			c.authenticated = true
			return nil
		}
	}
}

func (c *Client) makeJoins() {
	for _, streamer := range c.streamers {
		fmt.Fprintf(c, "JOIN %s\r\n", streamer)
	}
}

func (c *Client) SetStreamers(s ...string) error {
	if len(s) > 50 {
		return ErrTooMuchStreamers
	}
	c.streamers = s
	return nil
}

func (c *Client) AddStreamers(s ...string) error {
	if len(s)+len(c.streamers) > 50 {
		return ErrTooMuchStreamers
	}
	c.streamers = append(c.streamers, s...)
	return nil
}

func (c *Client) SendPong(rawPing string) {
	pong := strings.Replace(rawPing, "PING", "PONG", 1)
	fmt.Fprintf(c, "%s\r\n", pong)
}

func (c *Client) Write(b []byte) (n int, err error) {
	return c.conn.Write(b)
}
