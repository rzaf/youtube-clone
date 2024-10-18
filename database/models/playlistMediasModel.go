package models

import (
	"fmt"
	"log"
	"math"
	"time"
	"youtube-clone/database/db"
	"youtube-clone/database/pbs/helper"

	"github.com/lib/pq"
)

type PlaylistMedia struct {
	PlaylistUrl int64      `json:"playlist_url"`
	Text        string     `json:"text"`
	CreatedAt   *time.Time `json:"created_at"`
	MediaId     int64      `json:"-"`
	UserId      int64      `json:"-"`
	Order       int        `json:"order"`
	/// extra columns from other tables
	MediaType           helper.MediaType `json:"media_type"`
	MediaUrl            string           `json:"media_url"`
	MediaTitle          string           `json:"media_title"`
	MediaThumbnail      string           `json:"media_thumbnail"`
	MediaTime           *time.Time       `json:"media_created_at"`
	MediaCreator        string           `json:"media_creator"`
	MediaCreatorProfile string
	MediaViews          int64
}

//// Get

// all users can get medias from playlist
func GetMediasOfPlaylist(playlistUrl string, limit int, offset int) (int64, []*PlaylistMedia, error) {
	var totalPages int64
	var playlistId int64
	{
		query := `
		SELECT 
			playlist_id,COUNT(*)
		FROM 
			playlist_medias
		WHERE playlist_id=(SELECT getPlaylistIdByUrl($1))
		GROUP BY playlist_id;
		`
		rows, err := db.Db.Query(query, playlistUrl)
		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				log.Printf("%+v\n", *err)
				if err.Code == "P0001" {
					return 0, nil, &ModelError{Message: err.Message, Status: 404}
				}
			}
			return 0, nil, err
		}
		defer rows.Close()
		if !rows.Next() {
			return 0, nil, nil
		}
		err = rows.Scan(&playlistId, &totalPages)
		if err != nil {
			return 0, nil, err
		}
		fmt.Println(totalPages)
		totalPages = int64(math.Ceil(float64(totalPages) / float64(limit)))
	}

	query := `
	SELECT 
		M.url,
		M.title,
		COALESCE(M.thumbnail,'') AS thumbnail,
		M.media_type,
		(SELECT COUNT(*) FROM views WHERE views.media_id=M.id) AS views_count,
		U.username,
		COALESCE(U.profile_photo,'') AS creator_profile,
		M.created_at AS media_created_at,
		PM.created_at,
		PM.text,
		PM.custom_order
	FROM 
		playlist_medias PM
	JOIN
		medias M ON	
			M.id = PM.media_id
	JOIN
		users U ON
			U.id = M.user_id
	WHERE 
		playlist_id=$1 
	ORDER BY custom_order ASC ,created_at DESC 
	LIMIT $2 OFFSET $3;
	`
	rows, err := db.Db.Query(query, playlistId, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var medias []*PlaylistMedia
	for rows.Next() {
		var m PlaylistMedia
		err = rows.Scan(&m.MediaUrl, &m.MediaTitle, &m.MediaThumbnail, &m.MediaType, &m.MediaViews, &m.MediaCreator, &m.MediaCreatorProfile, &m.MediaTime, &m.CreatedAt, &m.Text, &m.Order)
		if err != nil {
			return 0, nil, err
		}
		medias = append(medias, &m)
	}
	// fmt.Printf("%+v\n", medias[0])

	return totalPages, medias, err
}

/// Create

// only creator of playlist should be able to add media to playlist
func AddMediaToPlaylist(playlistUrl string, mediaUrl string, text string, userId int64, customOrder int64) error {
	query := `
	INSERT INTO playlist_medias 
		(playlist_id,media_id,text,custom_order) 
	SELECT 
		P.id, 
		getMediaIdByUrlAndType($1,P.media_type),
		$2,
		$3
	FROM playlists P
	WHERE
		P.url=$4 AND P.user_id=$5 
	;
	`
	res, err := db.Db.Exec(query, mediaUrl, text, customOrder, playlistUrl, userId)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v\n", *err)
			if err.Code == "P0001" {
				return &ModelError{Message: err.Message, Status: 404}
			}
			if err.Code == "23505" {
				return NewModelError("media already exist in playlist!", 400)
			}
		}
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return NewModelError("playlist not exist or not belong to user!", 400)
	}
	return nil
}

/// Edit

// only creator of playlist should be able to edit media from playlist
func EditMediaFromPlaylist(newText string, newCustomOrder int64, playlistUrl string, mediaUrl string, userId int64) error {
	query := `
	UPDATE 
		playlist_medias 
	SET 
		text=$1,
		custom_order=$2,
		updated_at=$3 
	WHERE 
		media_id = getMediaIdByUrl($4)
		AND 
		playlist_id = (SELECT id FROM playlists WHERE playlists.url=$5 AND playlists.user_id=$6)
		--- playlist_id = getPlaylistIdByUrl($5)
		--- AND
		--- $6 = (SELECT user_id FROM playlists WHERE playlist_id=getPlaylistIdByUrl($2));
	`
	res, err := db.Db.Exec(query, newText, newCustomOrder, time.Now().UTC(), mediaUrl, playlistUrl, userId)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v\n", *err)
			if err.Code == "P0001" {
				return &ModelError{Message: err.Message, Status: 404}
			}
		}
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return NewModelError("media not exist in playlist or playlist not exist or not belong to user!", 400)
	}
	return nil
}

////// DELETE

// only creator of playlist should be able to remove media from playlist
func RemoveMediaFromPlaylist(playlistUrl string, mediaUrl string, userId int64) error {
	query := `
	DELETE FROM playlist_medias 
	WHERE 
		media_id = getMediaIdByUrl($1)
		AND 
		playlist_id = (SELECT id FROM playlists WHERE playlists.url=$2 AND playlists.user_id=$3)
		--- playlist_id = getPlaylistIdByUrl($2)
		--- AND
		--- $3 = (SELECT user_id FROM playlists WHERE playlist_id=getPlaylistIdByUrl($2));
	`
	res, err := db.Db.Exec(query, mediaUrl, playlistUrl, userId)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v\n", *err)
			if err.Code == "P0001" {
				return &ModelError{Message: err.Message, Status: 404}
			}
		}
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return NewModelError("media not exist in playlist or playlist not exist or not belong to user!", 400)
	}
	return nil
}
