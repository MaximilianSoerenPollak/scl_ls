package rpc_test

import (
	"sclls/rpc"
	"testing"
)

type EncodingExmpl struct {
	Testing bool
}

func TestEncodeMsg(t *testing.T) {
	expected := "Content-Length: 16\r\n\r\n{\"Testing\":true}"
	actual := rpc.EncodeMsg(EncodingExmpl{Testing: true})
	if expected != actual {
		t.Fatalf("Expected: %s, Actual: %s", expected, actual)
	}
}

func TestDecodeMsg(t *testing.T) {
	incMsg := "Content-Length: 15\r\n\r\n{\"Method\":\"hi\"}"
	// TODO Add content testing
	method, content, err := rpc.DecodeMsg([]byte(incMsg))
	contentLenght := len(content)
	if err != nil {
		t.Fatal(err)
	}
	if 15 != contentLenght {
		t.Fatalf("Expected: 15, Actual: %d", contentLenght)
	}
	if method != "hi" {
		t.Fatalf("Expected: 'hi', Actual: %s", method)
	}
}
