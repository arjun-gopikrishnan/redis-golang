// Note to self : circular import bw controller and peer. Check golang docs and best practices to avoid this scenario
package controller

import (
	"log"
	"log/slog"
	"strings"

	. "github.com/arjun/redis-go/internal/peer"
	"github.com/arjun/redis-go/internal/resp"
)

func CommandExecutor(parsedBuffer string,p *Peer){
	RedisCommand,err := resp.ParseRespCommand(parsedBuffer)

		if(err != nil){
			slog.Warn("Encountered an error while parsing text","Error" ,err.Error())
			response := resp.EncodeErrorMsgToResp(err.Error(),"ERR")
			p.WriteToClient(response)
			// continue
		}
		lines := strings.Split(strings.ReplaceAll(parsedBuffer, "\r", ""), ` `)
		action := strings.ToUpper(lines[0])
		switch action {
		case "PING":
			response := resp.EncodeSimpleStringToResp("OK")
			p.WriteToClient(response)
		case "ECHO":
			if len(RedisCommand.Args) > 0 {
				response := resp.EncodeBulkStringToResp(strings.Join(RedisCommand.Args," "))
				p.WriteToClient(response)
			} else {
				response := resp.EncodeNullToResp()
				p.WriteToClient(response)
			}
		case "SET":
			if len(RedisCommand.Args) < 2 {
				response := resp.EncodeErrorMsgToResp("SET needs key and value name","SYNTAXERR")
				p.WriteToClient(response)
				// continue
			}
			keyName := RedisCommand.Args[0]
			keyValue := strings.Join(RedisCommand.Args[1:], "\n")
			p.SetKey(keyName, keyValue)
			response := resp.EncodeSimpleStringToResp("OK")
			p.WriteToClient(response)
		case "GET":
			log.Println("Command called", "command", "GET")
			if len(lines) < 2 {
				response := resp.EncodeErrorMsgToResp("GET needs key name","SYNTAXERR")
				p.WriteToClient(response)
				// continue
			}

			keyName := RedisCommand.Args[0]

			keyValue, err := p.GetKey(keyName)

			if err != nil {
				response := resp.EncodeNullToResp()
				p.WriteToClient(response)
				// continue
			}
			response := resp.EncodeBulkStringToResp(keyValue)
			p.WriteToClient(response)
			
			

		default:
			p.WriteToClient(resp.EncodeErrorMsgToResp("Invalid command","ERR"))
		}
}