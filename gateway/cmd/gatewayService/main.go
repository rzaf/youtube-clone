package main

import (
	"log"

	"github.com/rzaf/youtube-clone/auth/middlewares"
	"github.com/rzaf/youtube-clone/gateway/client"
	"github.com/rzaf/youtube-clone/gateway/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	client.ConnectToDatabaseServer()
	defer client.DisconnectFromDatabaseServer()
	middlewares.SetSigningKey()
	server.StartHttpServer()
}
