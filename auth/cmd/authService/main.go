package main

import (
	"log"

	"github.com/rzaf/youtube-clone/auth/client"
	"github.com/rzaf/youtube-clone/auth/handlers"
	"github.com/rzaf/youtube-clone/auth/middlewares"
	"github.com/rzaf/youtube-clone/auth/server"
	"github.com/rzaf/youtube-clone/database/helpers"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	client.ConnectToDataBaseServer()
	defer client.DisconnectFromServer()

	handlers.SigningKey = []byte(helpers.FatalIfEmptyVar("JWT_SIGNING_KEY"))
	middlewares.SigningKey = []byte(helpers.FatalIfEmptyVar("JWT_SIGNING_KEY"))
	server.StartHttpServer()
}
