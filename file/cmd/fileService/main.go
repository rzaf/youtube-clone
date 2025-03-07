package main

import (
	"log"

	"github.com/rzaf/youtube-clone/auth/middlewares"
	"github.com/rzaf/youtube-clone/database/helpers"
	"github.com/rzaf/youtube-clone/file/client"
	"github.com/rzaf/youtube-clone/file/db"
	"github.com/rzaf/youtube-clone/file/queue"
	"github.com/rzaf/youtube-clone/file/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	client.ConnectToDataBaseServer()
	defer client.DisconnectFromServer()
	db.Connect()
	go queue.RunQueue()
	go server.StartGrpcServer()
	defer server.StopGrpcServer()

	middlewares.SigningKey = []byte(helpers.FatalIfEmptyVar("JWT_SIGNING_KEY"))
	server.StartHttpServer()
}
