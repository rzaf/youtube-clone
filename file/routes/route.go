package routes

import (
	"youtube-clone/file/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func GetRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)
	router.Use(handlers.RecoverServerPanics)

	router.Get("/photos/{url}", handlers.GetPhoto)
	router.Get("/videos/{url}", handlers.GetVideo)
	router.Get("/musics/{url}", handlers.GetMusic)

	router.Group(func(r chi.Router) {
		r.Use(handlers.AuthApiKeyMiddleware)
		r.Post("/photos/upload", handlers.UploadPhoto)
		r.Post("/videos/upload", handlers.UploadVideo)
		r.Post("/musics/upload", handlers.UploadMusic)
	})

	baseRouter := chi.NewRouter()
	baseRouter.Mount("/api", router)
	return baseRouter
}
