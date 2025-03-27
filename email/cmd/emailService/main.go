package main

import (
	"log"

	"github.com/rzaf/youtube-clone/email/email"
	"github.com/rzaf/youtube-clone/email/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	email.InitVars()
	defer server.StopGrpcServer()
	server.StartGrpcServer()
}
