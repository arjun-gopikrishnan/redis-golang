package resp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type RedisCommand struct{
	Command string
	Args []string
	Raw string
}

func (rc *RedisCommand) Display(){
	fmt.Print("Command given :" ,rc.Command)

	fmt.Print("Args : ")
	for _,s := range rc.Args{
		fmt.Print(s)
	}

	fmt.Print("Raw : ",rc.Raw)
}

func ParseRespCommand(s string) (RedisCommand,error){
	lines := strings.Split(strings.ReplaceAll(s, "\r\n", ""), ` `)
	if len(lines) < 1 {
		return RedisCommand{}, errors.New("empty command")
	}
	var command string
	switch strings.ToUpper(lines[0]) {
	case "PING", "SET", "GET", "ECHO":
		command = lines[0]
	default:
		return RedisCommand{}, errors.New("invalid command")
	}
	

	rObj := RedisCommand{
		Command: command,
		Args:    lines[1:],
		Raw:     s,
	}
	return rObj,nil
}

func EncodeNullToResp() string{
	return "$-1\r\n"
}

func EncodeErrorMsgToResp(errorMsg string,errPrompt string) string{
	return "- " + errPrompt + " " + errorMsg + "\r\n"
}

func EncodeBulkStringToResp(s string) string{
	lengthStr := strconv.Itoa(len(s))

	return "$" + lengthStr + "\r\n" + s + "\r\n"
}

func EncodeSimpleStringToResp(s string) string{
	return "+" + s + "\r\n"
}

func EncodeIntToResp(n int) string {
	var signBit string
	if n < 0 {
		signBit = "-"
	} else {
		signBit = "+"
	}
	return ":" + signBit + strconv.Itoa(n)
}

func EncodeErrorToResp(s string) string{
	return "-" + s + "\r\n"
}

