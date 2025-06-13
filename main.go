package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"sclls/lsp"
	"sclls/rpc"
)

func main() {
	logger := getLogger("/home/maxi/dev/scl_ls/log.txt")
	logger.Println("Hey, sclls started")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	for scanner.Scan() {
		msg := scanner.Bytes()
		method, content, err := rpc.DecodeMsg(msg)
		if err != nil {
			logger.Printf("got an error: %s", err.Error())
		}
		handleMessage(logger, method, content)
	}
}

func handleMessage(logger *log.Logger, method string, contents []byte) {
	logger.Printf("Revieced msg with method: %s", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("could not parse stuff: %s", err.Error())
		}
		logger.Printf("Connected to: %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)
		// let's reply here. How?
		writer := os.Stdout
		msg := lsp.NewInitializeReponse(request.ID)
		reply := rpc.EncodeMsg(msg)
		writer.Write([]byte(reply))

		logger.Printf("Send the reply: %v", msg)
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("could not parse stuff: %s", err.Error())
		}
		logger.Printf("Opened : %s", request.Params.TextDocument.URI)
		// let's reply here. How?
		logger.Printf("Text inside the File: %s", request.Params.TextDocument.Text)

	}
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		// TODO: better error here
		panic("not a good file")
	}
	return log.New(logfile, "[sclls]", log.Ldate|log.Ltime|log.Lshortfile)
}
