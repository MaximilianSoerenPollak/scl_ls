package rpc_test

import (
	"reflect"
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

func TestSplit(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		wantAdvance int
		wantToken   []byte
		wantErr     bool
	}{
		{
			name:        "complete valid message",
			input:       []byte("Content-Length: 13\r\n\r\nHello, World!"),
			wantAdvance: 35,
			wantToken:   []byte("Content-Length: 13\r\n\r\nHello, World!"),
			wantErr:     false,
		},
		{
			name:        "exact content length match",
			input:       []byte("Content-Length: 5\r\n\r\nHello"),
			wantAdvance: 26,
			wantToken:   []byte("Content-Length: 5\r\n\r\nHello"),
			wantErr:     false,
		},
		{
			name:        "incomplete message - no delimiter",
			input:       []byte("Content-Length: 13"),
			wantAdvance: 0,
			wantToken:   nil,
			wantErr:     false,
		},
		{
			name:        "invalid content length",
			input:       []byte("Content-Length: abc\r\n\r\nSome content"),
			wantAdvance: 0,
			wantToken:   nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAdvance, gotToken, err := rpc.Split(tt.input, false)

			// Check error expectation
			if tt.wantErr {
				if err == nil {
					t.Errorf("Split() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Split() unexpected error = %v", err)
				return
			}

			// Check advance
			if gotAdvance != tt.wantAdvance {
				t.Errorf("Split() advance = %v, want %v", gotAdvance, tt.wantAdvance)
			}

			// Check token
			if !reflect.DeepEqual(gotToken, tt.wantToken) {
				t.Errorf("Split() token = %v, want %v", gotToken, tt.wantToken)
				t.Errorf("Split() token string = %q, want %q", string(gotToken), string(tt.wantToken))
			}
		})
	}
}
