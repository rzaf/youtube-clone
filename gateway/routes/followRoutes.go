package routes

import (
	"youtube-clone/gateway/handlers"

	"github.com/go-chi/chi"
)

// func setFollowRoutes(router *chi.Mux) {
// }

func setFollowAuthRoutes(router chi.Router) {
	router.Post("/follows/{username}", handlers.AddFollowing)
	router.Delete("/follows/{username}", handlers.DeleteFollowing)
}
