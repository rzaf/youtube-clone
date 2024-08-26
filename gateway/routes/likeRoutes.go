package routes

import (
	"github.com/go-chi/chi"
	"youtube-clone/gateway/handlers"
)

// func setLikesRoutes(router *chi.Mux) {
// }

func setLikesAuthRoutes(router chi.Router) {
	router.Post("/medias/{url}/likes", handlers.SetMediaLike)
	router.Delete("/medias/{url}/likes", handlers.DeleteMediaLike)

	router.Post("/comments/{url}/likes", handlers.SetCommentLike)
	router.Delete("/comments/{url}/likes", handlers.DeleteCommentLike)
}
