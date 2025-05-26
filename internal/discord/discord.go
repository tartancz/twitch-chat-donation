package discord

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

var DefaultServer = NewServer()

type DiscordMessageArgs struct {
	Raw         string
	Args        []string
	CommandName string
}

type DiscordMessageHandler interface {
	HandleMessage(args DiscordMessageArgs, writer io.Writer)
	GetHelpMessage() string
}

type DiscordMessageHandlerStruct struct {
	HandleFunc func(args DiscordMessageArgs, writer io.Writer)
	HelpMessage string
}

func (f DiscordMessageHandlerStruct) HandleMessage(args DiscordMessageArgs, writer io.Writer) {
	if f.HandleFunc != nil {
		f.HandleFunc(args, writer)
	}
}

func (f DiscordMessageHandlerStruct) GetHelpMessage() string {
	if f.HelpMessage != "" {
		return f.HelpMessage
	}
	return "No help message provided."
}

type Server struct {
	messageChan chan []byte
	onMessage   func(string, io.Writer)
	handlers    map[string]DiscordMessageHandler
	conn        net.Conn
	delimiter   string
}

func NewServer() *Server {
	return &Server{
		messageChan: make(chan []byte),
	}
}

func (w *Server) Write(p []byte) (n int, err error) {
	if w.conn == nil {
		return 0, nil
	}

	msg := make([]byte, len(p), len(p)+len(w.delimiter)+1)
	copy(msg, p)

	msg = append(bytes.TrimSpace(msg), []byte("\n"+w.delimiter+"\n")...)

	w.messageChan <- msg

	return len(p), nil
}

func (w *Server) Close() {
	if w.messageChan != nil {
		close(w.messageChan)
		w.messageChan = nil
	}

}

func (w *Server) SetOnMessage(f func(string, io.Writer)) {
	w.onMessage = f
}

func (w *Server) AddHandler(command string, handler DiscordMessageHandler) {
	if w.handlers == nil {
		w.handlers = make(map[string]DiscordMessageHandler)
	}
	command = strings.ToLower(command)
	w.handlers[command] = handler
}

func (w *Server) HandleMessage(line string) {
	if w.onMessage != nil {
		w.onMessage(line, w)
	}

	line = strings.TrimSpace(line)
	args := strings.Split(strings.ToLower(line), " ")

	if len(args) < 3 || args[2] == "help" {
		fmt.Fprint(w, w.GeneretaHelp())
		return
	}
	command := args[2]
	if handler, exists := w.handlers[command]; exists {
		handler.HandleMessage(DiscordMessageArgs{
			Raw:  line,
			Args: args[3:],
			CommandName: command,
		}, w)
	} else {
		fmt.Fprintf(w, "Unknown command: %s\n", command)
		fmt.Fprint(w, w.GeneretaHelp())
	}
}
func (w *Server) GeneretaHelp() string {
	var helpMessage strings.Builder
	for command, handler := range w.handlers {
		helpMessage.WriteString(fmt.Sprintf("%s: %s\n-----------------\n", command, handler.GetHelpMessage()))
	}
	if helpMessage.Len() == 0 {
		return "No commands available."
	}
	return helpMessage.String()
}

func (w *Server) RunServer(programName string) {
	if w.conn != nil {
		return
	}

	port, exists := os.LookupEnv("DISCORD_BOT_SERVER_PORT")
	if !exists {
		return
	}
	host := getEnv("DISCORD_BOT_SERVER_HOST", "")
	var conn net.Conn
	var err error
	timeout := time.NewTimer(60 * time.Second)

connectLoop:
	for {
		select {
		case <-time.After(5 * time.Second):
			conn, err = net.Dial("tcp", net.JoinHostPort(host, port))
			if err != nil {
				fmt.Printf("Failed to connect to server: %v, retrying in 5 seconds...\n", err)
				continue
			}
			break connectLoop
		case <-timeout.C:
			fmt.Println("Connection timed out after 60 seconds")
			time.Sleep(60 * time.Minute)
			timeout.Reset(60 * time.Second)
		}
	}
	defer func() {
		conn.Close()
		w.conn = nil
	}()
	w.conn = conn
	fmt.Fprintf(conn, "SET_NAME:%s\n", programName)
	fmt.Fprintf(conn, "GET_DELIMITER:\n")

	go func() {
		reader := bufio.NewReader(conn)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			if strings.HasPrefix(line, "DELIMITER:") {
				w.delimiter = strings.TrimSpace(strings.TrimPrefix(line, "DELIMITER:"))
				continue
			}
			w.HandleMessage(line)
		}

	}()

	for val := range w.messageChan {
		conn.Write(val)
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
