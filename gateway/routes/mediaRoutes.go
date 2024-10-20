package routes

import (
	"github.com/rzaf/youtube-clone/gateway/handlers"

	"github.com/go-chi/chi"
)

func setMediaRoutes(router *chi.Mux) {
	router.Get("/medias", handlers.GetMedias)
	router.Get("/medias/search/{term}", handlers.SearchMedias)
	router.Get("/medias/{url}", handlers.GetMediaByUrl)
}

func setMediaAuthRoutes(router chi.Router) {

	router.Post("/medias", handlers.CreateMedia)
	router.Put("/medias/{url}", handlers.EditMedia)
	router.Delete("/medias/{url}", handlers.DeleteMedia)

	/// tag
	router.Post("/medias/{url}/tag/{name}", handlers.AddTagToVideo)
	router.Delete("/medias/{url}/tag/{name}", handlers.RemoveTagFromVideo)

	/// playlist
	router.Post("/medias/{url}/playlists/{playlistUrl}", handlers.AddMediaToPlaylist)
	router.Put("/medias/{url}/playlists/{playlistUrl}", handlers.EditMediaFromPlaylist)
	router.Delete("/medias/{url}/playlists/{playlistUrl}", handlers.DeleteMediaFromPlaylist)
}
