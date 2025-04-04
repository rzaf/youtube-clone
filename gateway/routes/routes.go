package routes

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rzaf/youtube-clone/gateway/handlers"

	authMiddleware "github.com/rzaf/youtube-clone/auth/middlewares"

	_ "github.com/rzaf/youtube-clone/gateway/docs"
	"github.com/swaggo/http-swagger" // http-swagger middleware
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
		r.Use(authMiddleware.JwtAuthMiddleware)

		setUserAuthRoutes(r)
		setMediaAuthRoutes(r)
		setFollowAuthRoutes(r)
		setLikesAuthRoutes(r)
		setCommentAuthRoutes(r)
		setPlaylistAuthRoutes(r)
	})

	baseRouter := chi.NewRouter()
	baseRouter.Mount("/api", router)
	baseRouter.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.UIConfig(map[string]string{
			"persistAuthorization": "true",
			"docExpansion":         "\"none\"",
		}),
	))
	return baseRouter
}
