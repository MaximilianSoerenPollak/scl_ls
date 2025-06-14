package main

import (
	"bufio"
	"encoding/json"
	"log"
	"flag"
	"os"

	"sclls/internal"
	"sclls/lsp"
	"sclls/rpc"

	//"github.com/goforj/godump"
)

func main() {
	logger := getLogger("/home/maxi/dev/scl_ls/log.txt")
	needsPath := flag.String("needsPath", "needs.json", "The path to your needs.json")
	enabled := flag.Bool("enable", true, "Disable the server.")
	docsPath := flag.String("docsPath", "docs", "The path to your docs folder")
	logger.Printf("Gotten following configs: %s, %s", needsPath, docsPath)
	logger.Println("Hey, sclls started")
	
	srvConfig := internal.ServerConfig{
		Enabled: *enabled,
		NeedsJsonPath: *needsPath,
		DocumentRootPath: *docsPath,
	}
	if !srvConfig.Enabled {
		logger.Println("Server was disabled. Exciting")
		os.Exit(0)
	}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	for scanner.Scan() {
		msg := scanner.Bytes()
		method, content, err := rpc.DecodeMsg(msg)
		if err != nil {
			logger.Printf("got an error: %s", err.Error())
		}
		handleMessage(logger, method, content, srvConfig)
	}
}

func handleMessage(logger *log.Logger, method string, contents []byte, srvConfig internal.ServerConfig) {
	logger.Printf("Revieced msg with method: %s", method)
	//logger.Printf("Revieced msg contents: %s", contents)

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
		documentNeedsEmpty := internal.NewDocumentNeeds(request.Params.TextDocument.URI, logger)
		content := []byte(request.Params.TextDocument.Text)
		ndi := internal.FindAllNeedsPositions(content, needsList)
		documentNeedsEmpty.Needs = ndi
		//out := godump.DumpStr(documentNeedsEmpty)
		//logger.Printf("FINISHED finding all needs: %v", out)

	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		logger.Printf("Revieced msg for did change contents: %s", contents)
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("could not parse stuff. didChange. Err: %s", err.Error())
		}
		logger.Printf("Opened : %s", request.Params.TextDocument.URI)
		for _, change := request.Params.ContentChanges {
			// TODO: Update text Document state	
		}
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
