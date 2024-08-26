package routes

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"youtube-clone/gateway/handlers"
)

func GetRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)
	router.Use(handlers.RecoverServerPanics)

	setMediaRoutes(router)
	setUserRoutes(router)
	setPlaylistRoutes(router)
	setCommentRoutes(router)

	router.Group(func(r chi.Router) {
		r.Use(handlers.AuthApiKeyMiddleware)

		setUserAuthRoutes(r)
		setMediaAuthRoutes(r)
		setFollowAuthRoutes(r)
		setLikesAuthRoutes(r)
		setCommentAuthRoutes(r)
		setPlaylistAuthRoutes(r)
	})

	baseRouter := chi.NewRouter()
	baseRouter.Mount("/api", router)
	return baseRouter
}
