package main

import (
	"log"

	"github.com/rzaf/youtube-clone/auth/middlewares"
	"github.com/rzaf/youtube-clone/notification/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	middlewares.SetSigningKey()
	go server.StartHttpServer()
	defer server.StopGrpcServer()
	server.StartGrpcServer()
}
