package main

import (
	"log"
	"youtube-clone/notification/email"
	"youtube-clone/notification/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	email.InitVars()
	defer server.StopGrpcServer()
	server.StartGrpcServer()
}
