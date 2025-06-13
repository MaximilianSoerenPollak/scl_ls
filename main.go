package main

import (
	"bufio"
	"os"
	"log"

	"sclls/rpc"

)

func main() {
	logger := getLogger("/home/maxi/dev/scl_ls/log.txt")
	logger.Println("Hey, sclls started")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	for scanner.Scan() {
		msg := scanner.Text()
		handleMessage(logger, msg)
	}
}

func handleMessage(logger *log.Logger, msg any) {
	logger.Println(msg)
}


func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		// TODO: better error here
		panic("not a good file")
	}
	return log.New(logfile, "[sclls]", log.Ldate|log.Ltime|log.Lshortfile)
}
