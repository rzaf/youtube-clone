package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	authMiddleware "github.com/rzaf/youtube-clone/auth/middlewares"
	"github.com/rzaf/youtube-clone/database/helpers"
	"github.com/rzaf/youtube-clone/notification/handlers"
)

func GetRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)
	router.Use(handlers.RecoverServerPanics)

	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.JwtAuthMiddleware)
		r.HandleFunc("/ws", handlers.WsHandler)
	})

	return router
}

func StartHttpServer() {
	host := helpers.FatalIfEmptyVar("WS_HOST")
	port := helpers.FatalIfEmptyVar("WS_PORT")

	log.Printf("listening at: %v:%v \n", host, port)
	log.Fatalln(http.ListenAndServe(host+":"+port, GetRoutes()))
}
