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

func EncodeArrayToResp(s string) string {
	stringArray := strings.Split(s, " ")
	wordCount := len(stringArray)
	var respArray strings.Builder

	// Append the array header
	respArray.WriteString(fmt.Sprintf("*%d\r\n", wordCount))

	// Encode each string and append to the result
	for _, substr := range stringArray {
		parsedString := EncodeBulkStringToResp(substr)
		respArray.WriteString(parsedString)
	}

	return respArray.String()
}

func DecodeArrayFromResp(s string) []string {
	RespCommands := strings.Split(s, "\r\n")
	var decodedArray []string

	// Ignore the first element and process the rest
	for i := 1; i < len(RespCommands); i += 2 {
		if i+1 < len(RespCommands) && RespCommands[i] != "" {
			decodedArray = append(decodedArray, RespCommands[i+1])
		}
	}

	return decodedArray
}

// func DecodeArrayToResp

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

