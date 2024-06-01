// tests/resp_test.go

package tests

import (
	"reflect"
	"testing"

	. "github.com/arjun/redis-go/internal/resp" // Adjust the import path based on your project structure
)

func TestParseRespCommand(t *testing.T) {
    // Define test cases
    testCases := []struct {
        input         string
        expectedResult RedisCommand
        expectedError error
    }{
        {
            input: "PING",
            expectedResult: RedisCommand{
                Command: "PING",
                Args:    []string{},
                Raw:     "PING",
            },
            expectedError: nil,
        },
        {
            input: "ECHO Hey There!",
            expectedResult: RedisCommand{
                Command: "ECHO",
                Args:    []string{"Hey","There!"},
                Raw:     "ECHO Hey There!",
            },
            expectedError: nil,
        },
        {
            input: "GET Key",
            expectedResult: RedisCommand{
                Command: "GET",
                Args:    []string{"Key"},
                Raw:     "GET Key",
            },
            expectedError: nil,
        },
        {
        input: "\r\nSET Key value\r\n",
        expectedResult: RedisCommand{
            Command: "SET",
            Args:    []string{"Key","value"},
            Raw:     "\r\nSET Key value\r\n",
        },
        expectedError: nil,
    },
        // Add more test cases as needed
    }

    // Run test cases
    for _, test := range testCases {
        result, err := ParseRespCommand(test.input)
        if !reflect.DeepEqual(result, test.expectedResult) {
            t.Errorf("Result mismatch for input: %s\nExpected: %+v\nActual: %+v", test.input, test.expectedResult, result)
        }

        if !reflect.DeepEqual(err, test.expectedError) {
            t.Errorf("Error mismatch for input: %s\nExpected: %+v\nActual: %+v", test.input, test.expectedError, err)
        }
    }
}


func TestEncodeBulkStringToResp(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"hey", "$3\r\nhey\r\n"},
		{"hello", "$5\r\nhello\r\n"},
		{"", "$0\r\n\r\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := EncodeBulkStringToResp(tc.input)
			println(result)
			if result != tc.expected {
				t.Errorf("Expected %s but got %s", tc.expected, result)
			}
		})
	}
}
