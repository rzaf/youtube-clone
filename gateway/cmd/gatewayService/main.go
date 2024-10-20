package main

import (
	"github.com/rzaf/youtube-clone/gateway/client"
	"github.com/rzaf/youtube-clone/gateway/server"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	client.ConnectToDatabaseServer()
	defer client.DisconnectFromDatabaseServer()
	server.StartHttpServer()
}
