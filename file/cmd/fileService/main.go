package main

import (
	"log"
	"youtube-clone/file/client"
	"youtube-clone/file/db"
	"youtube-clone/file/queue"
	"youtube-clone/file/server"
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
