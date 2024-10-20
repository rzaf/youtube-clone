package main

import (
	"github.com/rzaf/youtube-clone/notification/email"
	"github.com/rzaf/youtube-clone/notification/server"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	email.InitVars()
	defer server.StopGrpcServer()
	server.StartGrpcServer()
}
