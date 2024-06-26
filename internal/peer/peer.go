package peer

import (
	"bufio"
	"errors"
	"log/slog"
	"net"
	"strconv"
	"strings"

	"log"

	"github.com/arjun/redis-go/internal/keystore"
	. "github.com/arjun/redis-go/internal/resp"
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
		return "",  errors.New("empty command")
	}
	return value, nil
}

func (p *Peer) ReadLoop(s *keystore.Store) {
	defer p.conn.Close()
	scanner := bufio.NewScanner(p.conn)

	for scanner.Scan() {
		command := scanner.Text()
		RedisCommand,err := ParseRespCommand(command)

		if(err != nil){
			slog.Warn("Encountered an error while parsing text","Error" ,err.Error())
			response := EncodeErrorMsgToResp(err.Error(),"ERR")
			p.conn.Write([]byte(response))
			continue
		}
		lines := strings.Split(strings.ReplaceAll(command, "\r", ""), ` `)
		action := strings.ToUpper(lines[0])
		switch action {
		case "PING":
			response := EncodeSimpleStringToResp("OK")
			p.conn.Write([]byte(response))
		case "ECHO":
			if len(RedisCommand.Args) > 0 {
				response := EncodeBulkStringToResp(strings.Join(RedisCommand.Args," "))
				p.conn.Write([]byte(response))
			} else {
				response := EncodeNullToResp()
				p.conn.Write([]byte(response))
			}
		case "SET":
			
			var keyValue string
			lastCommandAt := len(RedisCommand.Args)
			
			ttl := 10000

			if len(RedisCommand.Args) < 2 {
				response := EncodeErrorMsgToResp("SET needs key and value name","SYNTAXERR")
				p.conn.Write([]byte(response))
				continue
			}

			keyName := RedisCommand.Args[0]


			if len(RedisCommand.Args) > 3{
				commandLen := len(RedisCommand.Args)

				//Expiry command handling
				if RedisCommand.Args[commandLen-2] == "px"{
					 
					ttl,err = strconv.Atoi(RedisCommand.Args[commandLen-1])

					if err!=nil{
						response := EncodeErrorMsgToResp("SET PX needs to be valid number","SYNTAXERR")
						p.conn.Write([]byte(response))
						continue
					}

					lastCommandAt = commandLen - 2   
				}
			}

			keyValue = strings.Join(RedisCommand.Args[1:lastCommandAt], "\n")
			p.SetKey(keyName, keyValue)
			s.SetKey(keyName,keyValue,ttl,"Arjun")
			response := EncodeSimpleStringToResp("OK")
			p.conn.Write([]byte(response))
		case "GET":
			log.Println("Command called", "command", "GET")
			if len(lines) < 2 {
				response := EncodeErrorMsgToResp("GET needs key name","SYNTAXERR")
				p.conn.Write([]byte(response))
				continue
			}

			keyName := RedisCommand.Args[0]

			// keyValue, err := p.GetKey(keyName)
			keyValue,err := s.GetKey(keyName)

			stringKeyVal := keyValue.Value()

			// fmt.Printf("%s",storeErr.Error())

			// fmt.Print(storeVal)
			if err != nil {
				response := EncodeNullToResp()
				p.conn.Write([]byte(response))
				continue
			}
			successResponse := EncodeBulkStringToResp(stringKeyVal)
			p.conn.Write([]byte(successResponse))
			
			

		default:
			p.conn.Write([]byte(EncodeErrorMsgToResp("Invalid command","ERR")))
		}

	}

}

func (p *Peer) WriteToClient(s string){
	p.conn.Write([]byte(s))
}