package routes

import (
	"github.com/rzaf/youtube-clone/file/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	authMiddleware "github.com/rzaf/youtube-clone/auth/middlewares"
	_ "github.com/rzaf/youtube-clone/file/docs"
	"github.com/swaggo/http-swagger" // http-swagger middleware
)

func GetRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)
	router.Use(handlers.RecoverServerPanics)

	router.Get("/photos/{url}", handlers.GetPhoto)
	router.Get("/videos/{url}", handlers.GetVideo)
	router.Get("/musics/{url}", handlers.GetMusic)

	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.JwtAuthMiddleware)
		r.Post("/photos/upload", handlers.UploadPhoto)
		r.Post("/videos/upload", handlers.UploadVideo)
		r.Post("/musics/upload", handlers.UploadMusic)
	})

	baseRouter := chi.NewRouter()
	baseRouter.Mount("/api", router)
	baseRouter.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.UIConfig(map[string]string{
			"persistAuthorization": "true",
			"docExpansion":         "\"none\"",
		}),
		httpSwagger.AfterScript(`document.title = "File Service";`),
	))
	return baseRouter
}
