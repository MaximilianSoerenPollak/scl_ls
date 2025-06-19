package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"log"
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
	//logger.Printf("Gotten following configs: %s, %s", needsPath, docsPath)
	logger.Println("Hey, sclls started")

	srvConfig := internal.ServerConfig{
		Enabled:          *enabled,
		NeedsJsonPath:    *needsPath,
		DocumentRootPath: *docsPath,
	}
	state := internal.NewState(srvConfig, logger)
	if !srvConfig.Enabled {
		logger.Println("Server was disabled. Exciting")
		os.Exit(0)
	}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	writer := os.Stdout
	for scanner.Scan() {
		msg := scanner.Bytes()
		method, content, err := rpc.DecodeMsg(msg)
		if err != nil {
			logger.Printf("got an error: %s", err.Error())
		}
		handleMessage(logger, writer, &state, method, content)
	}
}

func handleMessage(logger *log.Logger, writer io.Writer, state *internal.State, method string, contents []byte) {
	logger.Printf("Revieced msg with method: %s", method)
	//logger.Printf("Revieced msg contents: %s", contents)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("could not parse stuff: %s", err.Error())
			return
		}
		logger.Printf("Connected to: %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)
		// let's reply here. How?
		msg := lsp.NewInitializeReponse(request.ID)
		writeResponse(writer, msg)

		logger.Printf("Send the reply: %v", msg)
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("could not parse stuff: %s", err.Error())
		}
		logger.Printf("Opened : %s", request.Params.TextDocument.URI)
		// let's reply here. How?
		logger.Printf("Text inside the File: %s", request.Params.TextDocument.Text)
		state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)

	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("could not parse stuff. didChange. Err: %s", err.Error())
			return
		}
		logger.Printf("Opened : %s", request.Params.TextDocument.URI)
		for _, change := range request.Params.ContentChanges {
			state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
			// TODO: Update text Document state
		}
		// let's reply here. How?

	case "textDocument/hover":
		//Hover msg ('K')
		var request lsp.HoverRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("could not parse stuff: %s", err.Error())
			return
		}
		logger.Printf("Hover was requested")
		// Create a response
		// let's reply here. How?
		var responseStr string
		foundNeed, err := state.FindNeedsInRequestedPosition(request.Params.TextDocument.URI, request.Params.Position)
		if err != nil {
			responseStr = err.Error()
		} else {
			responseStr = foundNeed.GenerateHoverInfo()
		}
		response := lsp.HoverResponse{
			Response: lsp.Response{
				RPC: "2.0",
				ID:  &request.ID,
			},
			Result: lsp.HoverResult{
				Contents: responseStr,
			},
		}
		writeResponse(writer, response)
	case "textDocument/definition":
		//Def request ('K')
		var request lsp.DefinitionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("could not parse stuff in textDocument Definition request: %s", err.Error())
		}
		logger.Printf("Go to definition was requested")
		logger.Printf("ID: %d, URI: %s, Pos: %v", request.ID, request.Params.TextDocument.URI, request.Params.Position)

		// Create a response
		// let's reply here. How?
		msg := state.GoToDefinition(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(writer, msg)
	case "textDocument/completion":
		var request lsp.CompletionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("could not parse stuff in textDocument Completion request: %s", err.Error())
		}

		// Create a response
		// let's reply here. How?
		msg := state.TextDocumentCompletion(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(writer, msg)
	}
}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMsg(msg)
	writer.Write([]byte(reply))
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		// TODO: better error here
		panic("not a good file")
	}
	return log.New(logfile, "[sclls]", log.Ldate|log.Ltime|log.Lshortfile)
}
