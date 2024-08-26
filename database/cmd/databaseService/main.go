package main

import (
	"log"
	"youtube-clone/database/client"
	"youtube-clone/database/db"
	"youtube-clone/database/migrations"
	"youtube-clone/database/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	defer db.Disconnect()
	db.Connect()
	defer client.DisconnectFromFileServer()
	client.ConnectToFileServer()
	defer client.DisconnectFromNotificationServer()
	client.ConnectToNotificationServer()
	migrations.UpAll()
	defer server.StopGrpcServer()
	server.StartGrpcServer()
}
