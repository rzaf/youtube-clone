package routes

import (
	// "github.com/rzaf/youtube-clone/auth/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/rzaf/youtube-clone/auth/docs"
	"github.com/rzaf/youtube-clone/auth/handlers"
	"github.com/swaggo/http-swagger" // http-swagger middleware
)

func GetRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)
	// router.Use(middleware.Recoverer)
	router.Use(handlers.RecoverServerPanics)

	router.Post("/login", handlers.Login)
	router.Post("/refresh", handlers.Refresh)

	router.Post("/register", handlers.Register)

	// router.Get("/videos/{url}", handlers.GetVideo)
	// router.Get("/musics/{url}", handlers.GetMusic)

	// router.Group(func(r chi.Router) {
	// 	r.Use(handlers.AuthApiKeyMiddleware)
	// 	r.Post("/photos/upload", handlers.UploadPhoto)
	// 	r.Post("/videos/upload", handlers.UploadVideo)
	// 	r.Post("/musics/upload", handlers.UploadMusic)
	// })

	baseRouter := chi.NewRouter()
	baseRouter.Mount("/api", router)
	baseRouter.Get("/docs/*", httpSwagger.Handler())
	return baseRouter
}
