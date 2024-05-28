package main

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strings"
)

type Peer struct {
	conn  net.Conn
	store map[string]string
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		conn:  conn,
		store: make(map[string]string),
	}
}

func (p *Peer) SetKey(key string, value string) {
	p.store[key] = value
}

func (p *Peer) GetKey(key string) (string, error) {
	value, exists := p.store[key]
	print(p.store)
	if !exists {
		return "", errors.New("key not found")
	}
	return value, nil
}

func (p *Peer) readLoop() {
	defer p.conn.Close()
	scanner := bufio.NewScanner(p.conn)

	for scanner.Scan() {
		command := scanner.Text()
		lines := strings.Split(strings.ReplaceAll(command, "\r", ""), ` `)
		action := strings.ToUpper(lines[0])
		slog.Info("Printing the chosen command", "action", action)
		switch action {
		case "PING":
			p.conn.Write([]byte("+PONG\r\n"))
		case "ECHO":
			if len(lines) > 1 {
				response := strings.Join(lines[1:], "\n") + "\r\n"
				p.conn.Write([]byte(response))
			} else {
				p.conn.Write([]byte("\r\n"))
			}
		case "SET":
			slog.Info("Command called", "command", "SET")
			if len(lines) < 3 {
				slog.Info("Invalid SET command", "command", "SET")
				p.conn.Write([]byte("Invalid SET command\r\n"))
				return
			}
			keyName := lines[1]
			keyValue := strings.Join(lines[2:], "\n") + "\r\n"
			p.SetKey(keyName, keyValue)
			response := fmt.Sprintf("Set Key: %s\r\n", keyName)
			p.conn.Write([]byte(response))
		case "GET":
			slog.Info("Command called", "command", "GET")
			if len(lines) < 2 {
				slog.Info("Invalid GET command", "command", "GET")
				p.conn.Write([]byte("Invalid GET command\r\n"))
				return
			}

			keyName := lines[1]

			keyValue, err := p.GetKey(keyName)

			if err != nil {
				response := fmt.Sprintf("Error: %s\r\n", err)
				p.conn.Write([]byte(response))
				return
			}
			response := fmt.Sprintf("Value: %s\r\n", keyValue)
			p.conn.Write([]byte(response))

		default:
			p.conn.Write([]byte("Command not recognized\r\n"))
		}

	}

	// for _, line := range lines {
	// 	fmt.Print(lines)
	// 	command := strings.ToUpper(line)
	// 	switch command {
	// 	case "PING":
	// 		p.conn.Write([]byte("+PONG\r\n"))
	// 	case "ECHO":
	// 		p.conn.Write([]byte("ECHO\r\n"))
	// 	default:
	// 		p.conn.Write([]byte("Command not recognized\r\n"))
	// 	}

	// }
}

// func (s *Server) handleConnection(conn net.Conn) {
// 	defer conn.Close()

// 	scanner := bufio.NewScanner(conn)

// 	for scanner.Scan() {
// 		command := scanner.Text()
// 		lines := strings.Split(strings.ReplaceAll(command, "\r", ""), `\n`)
// 		fmt.Print(lines)
// 		for _, line := range lines {
// 			if strings.ToUpper(line) == "PING" {
// 				conn.Write([]byte("+PONG\r\n"))
// 			} else {
// 				conn.Write([]byte("Command not recognized\r\n"))
// 			}
// 		}
// 	}
// }
