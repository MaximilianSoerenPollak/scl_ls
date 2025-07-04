package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

func EncodeMsg(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		// TODO: Do not panic here
		panic(err)
	}
	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

type BaseMessage struct {
	Method string `json:"method"`
}

func DecodeMsg(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return "", nil, errors.New("did not find header")
	}
	// Content-Length: <number>
	contentLenghtBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLenghtBytes))
	if err != nil {
		return "", nil, err
	}

	var baseMsg BaseMessage
	if err := json.Unmarshal(content[:contentLength], &baseMsg); err != nil {
		return "", nil, err
	}
	return baseMsg.Method, content[:contentLength], nil
}

// type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)
func Split(data []byte, _ bool) (advance int, token []byte, err error) {
	header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return 0, nil, nil
	}
	// Content-Length: <number>
	contentLenghtBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLenghtBytes))
	if err != nil {
		// 'LS got no idea what to do here => error'
		return 0, nil, err
	}
	if len(content) < contentLength {
		return 0, nil, nil
	}
	totalLenght := len(header) + 4 + contentLength

	return totalLenght, data[:totalLenght], nil
}
