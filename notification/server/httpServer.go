package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	authMiddleware "github.com/rzaf/youtube-clone/auth/middlewares"
	_ "github.com/rzaf/youtube-clone/notification/docs"
	"github.com/rzaf/youtube-clone/notification/handlers"
	"github.com/rzaf/youtube-clone/notification/helpers"
	httpSwagger "github.com/swaggo/http-swagger"
)

func GetRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)
	router.Use(handlers.RecoverServerPanics)

	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.JwtAuthMiddleware)
		r.HandleFunc("/ws", handlers.WsHandler)

		r.HandleFunc("/api/notifications", handlers.GetNotifications)
		r.HandleFunc("/api/notifications/{id}", handlers.GetNotification)
		r.HandleFunc("/api/notifications/seen", handlers.ReadAllNotifications)
		r.HandleFunc("/api/notifications/{id}/seen", handlers.ReadNotification)
	})

	router.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.UIConfig(map[string]string{
			"persistAuthorization": "true",
			"docExpansion":         "\"none\"",
		}),
		httpSwagger.AfterScript(`document.title = "Notification Service";`),
	))
	return router
}

func StartHttpServer() {
	host := helpers.FatalIfEmptyVar("WS_HOST")
	port := helpers.FatalIfEmptyVar("WS_PORT")

	log.Printf("listening at: %v:%v \n", host, port)
	log.Fatalln(http.ListenAndServe(host+":"+port, GetRoutes()))
}
