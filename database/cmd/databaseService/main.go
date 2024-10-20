package main

import (
	"github.com/rzaf/youtube-clone/database/client"
	"github.com/rzaf/youtube-clone/database/db"
	"github.com/rzaf/youtube-clone/database/migrations"
	"github.com/rzaf/youtube-clone/database/server"
	"log"
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
