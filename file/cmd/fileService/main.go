package main

import (
	"github.com/rzaf/youtube-clone/file/client"
	"github.com/rzaf/youtube-clone/file/db"
	"github.com/rzaf/youtube-clone/file/queue"
	"github.com/rzaf/youtube-clone/file/server"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	client.ConnectToDataBaseServer()
	defer client.DisconnectFromServer()
	db.Connect()
	go queue.RunQueue()
	go server.StartGrpcServer()
	defer server.StopGrpcServer()

	server.StartHttpServer()
}
