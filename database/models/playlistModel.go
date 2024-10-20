package models

import (
	"errors"
	"fmt"
	"github.com/rzaf/youtube-clone/database/db"
	"github.com/rzaf/youtube-clone/database/helpers"
	"github.com/rzaf/youtube-clone/database/pbs/helper"
	"log"
	"math"
	"time"

	"github.com/lib/pq"
)

type Playlist struct {
	Id        int64      `json:"-"`
	Url       string     `json:"url"`
	Name      string     `json:"name"`
	Text      string     `json:"text"`
	Thumbnail string     `json:"thumbnail"`
	UserId    int64      `json:"-"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Type      helper.MediaType
	/// extra columns from other tables
	MediaCount      int64 `json:"medias_count"`
	MediaTotalViews int64 `json:"medias_views"`
	CreatorUsername string
	CreatorProfile  string
}

//// GetPlaylist

func GetPlaylist(playlistUrl string) (*Playlist, error) {
	query := `
	SELECT 
		P.name,
		P.text,
		P.media_type,
		P.created_at,
		P.updated_at,
		COALESCE(P.thumbnail,'') AS thumbnail,
		P.media_type,
		U.username,
		COALESCE(U.profile_photo,'') AS profile_photo,
		(SELECT COUNT(*) FROM playlist_medias PM WHERE PM.playlist_id=P.id) AS media_count
	FROM 
		playlists P
	JOIN users U
		ON P.user_id = U.id
	WHERE 
		url=$1;
	`
	rows, err := db.Db.Query(query, playlistUrl)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, NewModelError("playlist with url:`"+playlistUrl+"` not found", 404)
	}
	var p Playlist
	err = rows.Scan(&p.Name, &p.Text, &p.Type, &p.CreatedAt, &p.UpdatedAt, &p.Thumbnail, &p.Type, &p.CreatorUsername, &p.CreatorProfile, &p.MediaCount)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func SearchPlaylists(search string, limit int, offset int, sortType helper.SortType) (int64, []*Playlist, error) {
	var totalPages int64
	{
		query := `
		SELECT 
			COUNT(*)
		FROM 
			playlists P
		WHERE 
			P.name LIKE concat('%',$1::VARCHAR,'%')
		GROUP BY P.user_id;
		`
		rows, err := db.Db.Query(query, search)
		if err != nil {
			return 0, nil, err
		}
		defer rows.Close()
		if !rows.Next() {
			// return 0, nil, NewModelError("No playlist found", 404)
			return 0, nil, nil
		}
		err = rows.Scan(&totalPages)
		if err != nil {
			return 0, nil, err
		}
		totalPages = int64(math.Ceil(float64(totalPages) / float64(limit)))
	}

	st, err := helpers.SortToString(sortType)
	if err != nil {
		return 0, nil, err
	}
	query := `
	SELECT 
		P.name,
		P.url,
		P.text,
		P.created_at,
		COALESCE(P.thumbnail,'') AS thumbnail,
		U.username,
		(SELECT COUNT(*) FROM playlist_medias PM WHERE PM.playlist_id=P.id) AS media_count,
		(SELECT COUNT(*) FROM views V WHERE V.media_id IN (SELECT PM.media_id FROM playlist_medias PM WHERE PM.playlist_id=P.id)) AS views_count
	FROM 
		playlists P
	JOIN users U
		ON U.id=P.user_id
	WHERE
		P.name LIKE concat('%',$1::VARCHAR,'%')
	ORDER BY ` + st + `
	LIMIT $2 OFFSET $3;	
	`
	fmt.Println(limit, offset)
	rows, err := db.Db.Query(query, search, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var playlists []*Playlist
	for rows.Next() {
		var p Playlist
		err = rows.Scan(&p.Name, &p.Url, &p.Text, &p.CreatedAt, &p.Thumbnail, &p.CreatorUsername, &p.MediaCount, &p.MediaTotalViews)
		if err != nil {
			return 0, nil, err
		}
		playlists = append(playlists, &p)
	}
	return totalPages, playlists, nil
}

func SearchUserPlaylists(search string, userName string, limit int, offset int, sortType helper.SortType) (int64, []*Playlist, error) {
	var totalPages int64
	var userId int64
	{
		query := `
		SELECT 
			P.user_id,COUNT(*)
		FROM 
			playlists P
		WHERE 
			P.user_id=getUserIdByUsername($1) AND P.name LIKE concat('%',$2::VARCHAR,'%')
		GROUP BY P.user_id;
		`
		rows, err := db.Db.Query(query, userName, search)
		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				log.Printf("%+v\n", *err)
				if err.Code == "P0001" { // exception
					return 0, nil, &ModelError{Message: err.Message, Status: 400}
				}
			}
			return 0, nil, err
		}
		defer rows.Close()
		if !rows.Next() {
			// return 0, nil, NewModelError("No playlist found", 404)
			return 0, nil, nil
		}
		err = rows.Scan(&userId, &totalPages)
		if err != nil {
			return 0, nil, err
		}
		totalPages = int64(math.Ceil(float64(totalPages) / float64(limit)))
	}

	st, err := helpers.SortToString(sortType)
	if err != nil {
		return 0, nil, err
	}
	query := `
	SELECT 
		P.name,
		P.url,
		P.text,
		P.created_at,
		COALESCE(P.thumbnail,'') AS thumbnail,
		U.username,
		(SELECT COUNT(*) FROM playlist_medias PM WHERE PM.playlist_id=P.id) AS media_count,
		(SELECT COUNT(*) FROM views V WHERE V.media_id IN (SELECT PM.media_id FROM playlist_medias PM WHERE PM.playlist_id=P.id)) AS views_count
	FROM 
		playlists P
	JOIN users U
		ON U.id=P.user_id
	WHERE
		P.user_id=$1 AND P.name LIKE concat('%',$2::VARCHAR,'%')
	ORDER BY ` + st + `
	LIMIT $3 OFFSET $4;
	`
	fmt.Println(limit, offset)
	rows, err := db.Db.Query(query, userId, search, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var playlists []*Playlist
	for rows.Next() {
		var p Playlist
		err = rows.Scan(&p.Name, &p.Url, &p.Text, &p.CreatedAt, &p.Thumbnail, &p.CreatorUsername, &p.MediaCount, &p.MediaTotalViews)
		if err != nil {
			return 0, nil, err
		}
		playlists = append(playlists, &p)
	}
	return totalPages, playlists, nil
}

func GetPlaylists(limit int, offset int, sortType helper.SortType) (int64, []*Playlist, error) {
	var totalPages int64
	{
		query := `
		SELECT 
			COUNT(*)
		FROM 
			playlists;
		`
		rows, err := db.Db.Query(query)
		if err != nil {
			return 0, nil, err
		}
		defer rows.Close()
		if !rows.Next() {
			return 0, nil, nil
		}
		err = rows.Scan(&totalPages)
		if err != nil {
			return 0, nil, err
		}
		totalPages = int64(math.Ceil(float64(totalPages) / float64(limit)))
	}

	st, err := helpers.SortToString(sortType)
	if err != nil {
		return 0, nil, err
	}
	query := `
	SELECT 
		P.name,
		P.url,
		P.text,
		P.created_at,
		COALESCE(P.thumbnail,'') AS thumbnail,
		U.username,
		(SELECT COUNT(*) FROM playlist_medias PM WHERE PM.playlist_id=P.id) AS media_count,
		(SELECT COUNT(*) FROM views V WHERE V.media_id IN (SELECT PM.media_id FROM playlist_medias PM WHERE PM.playlist_id=P.id)) AS views_count
	FROM 
		playlists P
	JOIN users U
		ON U.id=P.user_id
	ORDER BY ` + st + `
	LIMIT $1 OFFSET $2;	
	`
	fmt.Println(limit, offset)
	rows, err := db.Db.Query(query, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var playlists []*Playlist
	for rows.Next() {
		var p Playlist
		err = rows.Scan(&p.Name, &p.Url, &p.Text, &p.CreatedAt, &p.Thumbnail, &p.CreatorUsername, &p.MediaCount, &p.MediaTotalViews)
		if err != nil {
			return 0, nil, err
		}
		playlists = append(playlists, &p)
	}
	return totalPages, playlists, nil
}

func GetUserPlaylists(userName string, limit int, offset int, sortType helper.SortType) (int64, []*Playlist, error) {
	var totalPages int64
	var userId int64
	{
		query := `
		SELECT 
			user_id,COUNT(*)
		FROM 
			playlists
		WHERE user_id=getUserIdByUsername($1)
		GROUP BY user_id;
		`
		rows, err := db.Db.Query(query, userName)
		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				log.Printf("%+v\n", *err)
				if err.Code == "P0001" { // exception
					return 0, nil, &ModelError{Message: err.Message, Status: 400}
				}
			}
			return 0, nil, err
		}
		defer rows.Close()
		if !rows.Next() {
			return 0, nil, nil
		}
		err = rows.Scan(&userId, &totalPages)
		if err != nil {
			return 0, nil, err
		}
		totalPages = int64(math.Ceil(float64(totalPages) / float64(limit)))
	}

	st, err := helpers.SortToString(sortType)
	if err != nil {
		return 0, nil, err
	}
	query := `
	SELECT 
		P.name,
		P.url,
		P.text,
		P.created_at,
		COALESCE(P.thumbnail,'') AS thumbnail,
		U.username,
		(SELECT COUNT(*) FROM playlist_medias PM WHERE PM.playlist_id=P.id) AS media_count,
		(SELECT COUNT(*) FROM views V WHERE V.media_id IN (SELECT PM.media_id FROM playlist_medias PM WHERE PM.playlist_id=P.id)) AS views_count
	FROM 
		playlists P
	JOIN users U
		ON U.id=P.user_id
	WHERE 
		P.user_id=$1
	ORDER BY ` + st + `
	LIMIT $2 OFFSET $3;	
	`
	rows, err := db.Db.Query(query, userId, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var playlists []*Playlist
	for rows.Next() {
		var p Playlist
		err = rows.Scan(&p.Name, &p.Url, &p.Text, &p.CreatedAt, &p.Thumbnail, &p.CreatorUsername, &p.MediaCount, &p.MediaTotalViews)
		if err != nil {
			return 0, nil, err
		}
		playlists = append(playlists, &p)
	}
	return totalPages, playlists, nil
}

/// Create

func CreatePlaylist(name string, text string, mt helper.MediaType, userId int64) (string, error) {
	query := "INSERT INTO playlists (name,text,media_type,user_id,url) VALUES ($1,$2,$3,$4,$5)"
	url := generateSecureToken(16)
	res, err := db.Db.Exec(query, name, text, mt, userId, url)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v\n", *err)
			// if err.Code == "P0001" { // exception
			// 	return "", &ModelError{Message: err.Message, Status: 400}
			// }
		}
		return "", err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return "", err
	}
	if n == 0 {
		return "", errors.New("failed to create playlist")
	}
	return url, nil
}

/// Edit

func EditPlaylist(name string, text string, playlistUrl string, userId int64) error {
	query := "UPDATE playlists SET name=$1,text=$2,updated_at=$3 WHERE user_id=$4 AND id=getPlaylistIdByUrl($5) ;"
	res, err := db.Db.Exec(query, name, text, time.Now().UTC(), userId, playlistUrl)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v\n", *err)
			if err.Code == "P0001" { // exception
				return &ModelError{Message: err.Message, Status: 400}
			}
		}
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return NewModelError("playlist is not created by user", 403)
	}
	return nil
}

//// DELETE

func DeletePlaylist(playlistUrl string, userId int64) error {
	fmt.Println(userId, playlistUrl)
	query := "DELETE FROM playlists WHERE user_id=$1 AND id=(SELECT getPlaylistIdByUrl($2)) ;"
	res, err := db.Db.Exec(query, userId, playlistUrl)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v\n", *err)
			if err.Code == "P0001" { // exception
				return &ModelError{Message: err.Message, Status: 400}
			}
		}
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return NewModelError("playlist is not created by user", 403)
	}
	return nil
}
