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

	defer client.DisconnectFromNotificationServer()
	client.ConnectToNotificationServer()
	defer client.DisconnectFromFileServer()
	client.ConnectToFileServer()
	defer client.DisconnectFromEmailServer()
	client.ConnectToEmailServer()
	migrations.UpAll()
	defer server.StopGrpcServer()
	server.StartGrpcServer()
}
