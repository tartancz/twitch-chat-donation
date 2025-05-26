package twitch

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type TwitchAuth struct {
	Oauth string
	Nick  string
}

type TwitchReader struct {
	Auth    TwitchAuth
	Channel string
	Conn    net.Conn
	Reader  *bufio.Reader
}

func NewTwitchReader(auth TwitchAuth, channel string) (*TwitchReader, error) {
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Twitch IRC: %w", err)
	}

	// Set connection timeout
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	fmt.Fprintf(conn, "PASS %s\r\n", auth.Oauth)
	fmt.Fprintf(conn, "NICK %s\r\n", auth.Nick)
	fmt.Fprintf(conn, "JOIN %s\r\n", channel)

	fmt.Fprintf(os.Stdout, "PASS %s\r\n", auth.Oauth)
	fmt.Fprintf(os.Stdout, "NICK %s\r\n", auth.Nick)
	fmt.Fprintf(os.Stdout, "JOIN %s\r\n", channel)

	// Remove deadline after initial setup
	conn.SetDeadline(time.Time{})

	t := &TwitchReader{
		Auth:    auth,
		Channel: channel,
		Conn:    conn,
		Reader:  bufio.NewReader(conn),
	}
	return t, nil
}

func (t *TwitchReader) ReadLine() (string, error) {
	for {
		line, err := t.Reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read from connection: %w", err)
		}

		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "PING") {
			pong := strings.Replace(line, "PING", "PONG", 1)
			fmt.Fprintf(t.Conn, "%s\r\n", pong)
			continue
		}

		if !strings.Contains(line, "PRIVMSG") {
			continue
		}

		return line, nil
	}
}

func (t *TwitchReader) Close() error {
	if t.Conn != nil {
		return t.Conn.Close()
	}
	return nil
}

// Context-aware reading
func (t *TwitchReader) ReadLineWithContext(ctx context.Context) (string, error) {
	type result struct {
		line string
		err  error
	}

	resultChan := make(chan result, 1)
	go func() {
		line, err := t.ReadLine()
		resultChan <- result{line: line, err: err}
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case res := <-resultChan:
		return res.line, res.err
	}
}