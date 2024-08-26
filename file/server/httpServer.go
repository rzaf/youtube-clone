package server

import (
	"fmt"
	"log"
	"net/http"
	"youtube-clone/database/helpers"
	"youtube-clone/file/routes"
)

func StartHttpServer() {
	host := helpers.FatalIfEmptyVar("HTTP_HOST")
	port := helpers.FatalIfEmptyVar("HTTP_PORT")
	baseRouter := routes.GetRoutes()
	fmt.Printf("listening at: %v:%v \n", host, port)
	err := http.ListenAndServe(host+":"+port, baseRouter)
	if err != nil {
		log.Fatalln(err)
	}
}
