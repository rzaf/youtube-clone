package server

import (
	"fmt"
	"github.com/rzaf/youtube-clone/auth/routes"
	"github.com/rzaf/youtube-clone/database/helpers"
	"log"
	"net/http"
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
