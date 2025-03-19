package main

import (
	"log"

	"github.com/rzaf/youtube-clone/auth/client"
	"github.com/rzaf/youtube-clone/auth/handlers"
	"github.com/rzaf/youtube-clone/auth/middlewares"
	"github.com/rzaf/youtube-clone/auth/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	client.ConnectToDataBaseServer()
	defer client.DisconnectFromServer()

	middlewares.SetSigningKey()
	handlers.SigningKey = middlewares.GetSigningKey()
	server.StartHttpServer()
}
