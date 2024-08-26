package services

import (
	"context"
	"fmt"
	"log"
	"youtube-clone/database/models"
	"youtube-clone/database/pbs/helper"
	"youtube-clone/database/pbs/playlist"
)

type playlistServiceServer struct {
	playlist.PlaylistServiceServer
}

func newResponseFromPlaylistData(p *playlist.PlaylistData) *playlist.Response {
	return &playlist.Response{
		Res: &playlist.Response_Playlist{
			Playlist: p,
		},
	}
}

func newResponseFromPlaylists(playlists []*models.Playlist, pageInfo *helper.PagesInfo) *playlist.Response {
	if playlists == nil {
		return newPlaylistResponseFromEmpty()
	}
	var playlistsData []*playlist.PlaylistShortData
	for _, p := range playlists {
		p2 := &playlist.PlaylistShortData{
			Name:             p.Name,
			Url:              p.Url,
			CreatorUsername:  p.CreatorUsername,
			CreatorProfile:   p.CreatorProfile,
			Thumbnail:        p.Thumbnail,
			MediaCount:       p.MediaCount,
			MediasTotalViews: p.MediaTotalViews,
			CreatedAt:        p.CreatedAt.Unix(),
		}
		playlistsData = append(playlistsData, p2)
	}
	return &playlist.Response{
		Res: &playlist.Response_Playlists{
			Playlists: &playlist.PlaylistsData{
				Playlists: playlistsData,
				PagesInfo: pageInfo,
			},
		},
	}
}

func newPlaylistResponseFromError(e *helper.HttpError) *playlist.Response {
	return &playlist.Response{
		Res: &playlist.Response_Err{
			Err: e,
		},
	}
}

func newPlaylistResponseFromEmpty() *playlist.Response {
	return &playlist.Response{
		Res: &playlist.Response_Empty{},
	}
}

func (*playlistServiceServer) GetPlaylist(c context.Context, p *playlist.PlaylistReq) (*playlist.Response, error) {
	playlistM, err := models.GetPlaylist(p.PlaylistUrl)
	fmt.Println(playlistM)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newPlaylistResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	p2 := &playlist.PlaylistData{
		Url:              playlistM.Url,
		Name:             playlistM.Name,
		Text:             playlistM.Text,
		CreatedAt:        playlistM.CreatedAt.Unix(),
		CreatorUsername:  playlistM.CreatorUsername,
		CreatorProfile:   playlistM.CreatorProfile,
		Thumbnail:        playlistM.Thumbnail,
		MediaCount:       playlistM.MediaCount,
		MediaType:        mediaTypeToStr(playlistM.Type),
		MediasTotalViews: playlistM.MediaTotalViews,
	}
	if playlistM.UpdatedAt != nil {
		p2.UpdatedAt = playlistM.UpdatedAt.Unix()
	}
	return newResponseFromPlaylistData(p2), nil
}

func (*playlistServiceServer) GetPlaylists(c context.Context, p *playlist.PlaylistReq) (*playlist.Response, error) {
	PerPage, PageNumber := getPage(p.Page)
	var err error
	var totalPages int64
	var playlists []*models.Playlist
	if p.Username == "" {
		totalPages, playlists, err = models.GetPlaylists(PerPage, PageNumber, p.Sort)
	} else {
		totalPages, playlists, err = models.GetUserPlaylists(p.Username, PerPage, PageNumber, p.Sort)
	}
	fmt.Println(totalPages, playlists)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newPlaylistResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newResponseFromPlaylists(playlists, newPagesInfo(int32(totalPages), p.Page.PageNumber)), nil
}

func (*playlistServiceServer) SearchPlaylists(c context.Context, p *playlist.PlaylistReq) (*playlist.Response, error) {
	PerPage, PageNumber := getPage(p.Page)
	var err error
	var totalPages int64
	var playlists []*models.Playlist

	if p.Username == "" {
		totalPages, playlists, err = models.SearchPlaylists(p.SearchTerm, PerPage, PageNumber, p.Sort)
	} else {
		totalPages, playlists, err = models.SearchUserPlaylists(p.SearchTerm, p.Username, PerPage, PageNumber, p.Sort)
	}
	fmt.Println(totalPages, playlists)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newPlaylistResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newResponseFromPlaylists(playlists, newPagesInfo(int32(totalPages), p.Page.PageNumber)), nil
}

func (*playlistServiceServer) CreatePlaylist(c context.Context, p *playlist.EditPlaylistData) (*playlist.Response, error) {
	url, err := models.CreatePlaylist(p.Name, p.Text, helper.MediaType(p.MediaTypeId), p.CurrentUserId)
	log.Println(err)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newPlaylistResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newResponseFromPlaylistData(&playlist.PlaylistData{
		Url: url,
	}), nil
}

func (*playlistServiceServer) EditPlaylist(c context.Context, p *playlist.EditPlaylistData) (*playlist.Response, error) {
	err := models.EditPlaylist(p.Name, p.Text, p.Url, p.CurrentUserId)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newPlaylistResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newPlaylistResponseFromEmpty(), nil
}

func (*playlistServiceServer) DeletePlaylist(c context.Context, p *playlist.EditPlaylistData) (*playlist.Response, error) {
	err := models.DeletePlaylist(p.Url, p.CurrentUserId)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newPlaylistResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newPlaylistResponseFromEmpty(), nil
}

////////playlist medias

func newResponseFromPlaylistMedias(mediasM []*models.PlaylistMedia, pageInfo *helper.PagesInfo) *playlist.Response {
	if mediasM == nil {
		return newPlaylistResponseFromEmpty()
	}
	var medias []*playlist.Media
	for _, m := range mediasM {
		m2 := &playlist.Media{
			Note:                 m.Text,
			Order:                int32(m.Order),
			CreatedAt:            m.CreatedAt.Unix(),
			MediaTitle:           m.MediaTitle,
			MediaUrl:             m.MediaUrl,
			MediaCreatedAt:       m.MediaTime.Unix(),
			MediaThumbnail:       m.MediaThumbnail,
			MediaType:            mediaTypeToStr(m.MediaType),
			MediaCreatorUsername: m.MediaCreator,
			MediaCreatorProfile:  m.MediaCreatorProfile,
			MediaViews:           m.MediaViews,
		}
		medias = append(medias, m2)
	}
	return &playlist.Response{
		Res: &playlist.Response_Medias{
			Medias: &playlist.PlaylistMediasData{
				PlaylistMedias: medias,
				PagesInfo:      pageInfo,
			},
		},
	}
}
func (*playlistServiceServer) GetPlaylistMedias(c context.Context, p *playlist.PlaylistReq) (*playlist.Response, error) {
	PerPage, PageNumber := getPage(p.Page)
	totalPages, medias, err := models.GetMediasOfPlaylist(p.PlaylistUrl, PerPage, PageNumber)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newPlaylistResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newResponseFromPlaylistMedias(medias, newPagesInfo(int32(totalPages), p.Page.PageNumber)), nil
}

func (*playlistServiceServer) AddMediaToPlaylist(c context.Context, p *playlist.PlaylistMediaReq) (*playlist.Response, error) {
	err := models.AddMediaToPlaylist(p.PlaylistUrl, p.MediaUrl, p.Note, p.UserId, p.Order)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newPlaylistResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newPlaylistResponseFromEmpty(), nil
}

func (*playlistServiceServer) EditMediaFromPlaylist(c context.Context, p *playlist.PlaylistMediaReq) (*playlist.Response, error) {
	err := models.EditMediaFromPlaylist(p.Note, p.Order, p.PlaylistUrl, p.MediaUrl, p.UserId)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newPlaylistResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newPlaylistResponseFromEmpty(), nil
}

func (*playlistServiceServer) RemoveMediaFromPlaylist(c context.Context, p *playlist.PlaylistMediaReq) (*playlist.Response, error) {
	err := models.RemoveMediaFromPlaylist(p.PlaylistUrl, p.MediaUrl, p.UserId)
	if err != nil {
		if err2, ok := models.ConvertError(err); ok {
			return newPlaylistResponseFromError(err2.ToHttpError()), nil
		}
		return nil, err
	}
	return newPlaylistResponseFromEmpty(), nil
}
