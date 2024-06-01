package peer

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"log"
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
	if !exists {
		return "$-1\r\n", nil
	}
	return value, nil
}

func (p *Peer) ReadLoop() {
	defer p.conn.Close()
	scanner := bufio.NewScanner(p.conn)

	for scanner.Scan() {
		command := scanner.Text()
		lines := strings.Split(strings.ReplaceAll(command, "\r", ""), ` `)
		action := strings.ToUpper(lines[0])
		log.Println("Printing the chosen command", "action", action)
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
			log.Println("Command called", "command", "SET")
			if len(lines) < 3 {
				log.Println("Invalid SET command", "command", "SET")
				p.conn.Write([]byte("Invalid SET command\r\n"))
				return
			}
			keyName := lines[1]
			keyValue := strings.Join(lines[2:], "\n") + "\r\n"
			p.SetKey(keyName, keyValue)
			response := "+OK\r\n"
			p.conn.Write([]byte(response))
		case "GET":
			log.Println("Command called", "command", "GET")
			if len(lines) < 2 {
				log.Println("Invalid GET command", "command", "GET")
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
			response := keyValue
			p.conn.Write([]byte(response))

		default:
			p.conn.Write([]byte("Command not recognized\r\n"))
		}

	}

}
