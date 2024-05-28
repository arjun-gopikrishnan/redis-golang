package main

import (
	"bufio"
	"fmt"
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

	// if err != nil {
	// 	slog.Error("Error during buffer read", "err", err)
	// }

	for scanner.Scan() {
		command := scanner.Text()
		lines := strings.Split(strings.ReplaceAll(command, "\r", ""), `\n`)
		fmt.Print(lines)
		for _, line := range lines {
			if strings.ToUpper(line) == "PING" {
				p.conn.Write([]byte("+PONG\r\n"))
			} else {
				p.conn.Write([]byte("Command not recognized\r\n"))
			}
		}
	}
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
