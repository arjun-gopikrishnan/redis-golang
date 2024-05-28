package main

import (
	"bufio"
	"log/slog"
	"net"
	"strings"
)

type Peer struct {
	conn net.Conn
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		conn: conn,
	}
}

func (p *Peer) readLoop() {
	defer p.conn.Close()
	scanner := bufio.NewScanner(p.conn)

	for scanner.Scan() {
		command := scanner.Text()
		lines := strings.Split(strings.ReplaceAll(command, "\r", ""), `\n`)
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
