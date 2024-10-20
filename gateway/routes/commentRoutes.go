package routes

import (
	"github.com/rzaf/youtube-clone/gateway/handlers"

	"github.com/go-chi/chi"
)

func setCommentRoutes(router *chi.Mux) {
	router.Get("/comments/{commentUrl}", handlers.GetComment)

	router.Get("/comments/medias/{url}", handlers.GetCommentsOfMedia)
	router.Get("/comments/{commentUrl}/replies", handlers.GetRepliesOfComment)
}

func setCommentAuthRoutes(router chi.Router) {
	router.Post("/comments/medias/{url}", handlers.CreateComment)
	router.Put("/comments/{commentUrl}", handlers.EditComment)
	router.Delete("/comments/{commentUrl}", handlers.DeleteComment)
}
