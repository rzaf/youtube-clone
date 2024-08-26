package main

import (
	"log"
	"youtube-clone/gateway/client"
	"youtube-clone/gateway/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	client.ConnectToDatabaseServer()
	defer client.DisconnectFromDatabaseServer()
	server.StartHttpServer()
}
