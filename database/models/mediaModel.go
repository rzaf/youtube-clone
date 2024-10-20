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

// type MediaType int

// func (t MediaType) String() string {
// 	switch t {
// 	case 0:
// 		return "video"
// 	case 1:
// 		return "music"
// 	case 2:
// 		return "photo"
// 	}
// 	return ""
// }

// const (
// 	VIDEO MediaType = iota
// 	MUSIC
// 	PHOTO
// 	ALL
// )

type Media struct {
	Id        int64            `json:"id"`
	Title     string           `json:"title"`
	Text      string           `json:"text"`
	Url       string           `json:"url"`
	Type      helper.MediaType `json:"type"` //0:video 1:music 2:photo
	Thumbnail string           `json:"thumbnail"`
	UserId    int64            `json:"-"`
	CreatedAt *time.Time       `json:"created_at"`
	UpdatedAt *time.Time       `json:"updated_at,omitempty"`
	/// extra columns from other tables
	UserName                     string    `json:"uploader_username"`
	UserProfile                  string    `json:"uploader_profile"`
	Tags                         string    `json:"tags"` // seperated with comma
	ChannelName                  string    `json:"uploader_name"`
	UserSubscribesCount          int64     `json:"uploader_subscribers"`
	ViewsCount                   int64     `json:"views_count"`
	LikesCount                   int64     `json:"likes_count"`
	DislikesCount                int64     `json:"dislikes_count"`
	CommentsCount                int64     `json:"comments_count"`
	CurrentUserLike              LikeState `json:"-"`
	IsCurrentUserSubbedToCreator bool      `json:"-"`
}

func SearchMedias(limit int, offset int, search string, t helper.MediaType, sortType helper.SortType) (int64, []*Media, error) {
	mt := fmt.Sprintf(" M.media_type=%d AND ", t)
	if t == helper.MediaType_ALL {
		mt = ""
	}
	var totalPages int64
	{
		query := `
		SELECT 
			COUNT(*)
		FROM 
			medias AS M
		WHERE ` + mt + `M.title LIKE concat('%',$1::VARCHAR,'%');
		`
		rows, err := db.Db.Query(query, search)
		if err != nil {
			return 0, nil, err
		}
		defer rows.Close()
		if !rows.Next() {
			// return 0, nil, NewModelError("No media found", 404)
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
	fmt.Println(st)
	query := `
	SELECT 
		M.title,
		M.url,
		M.media_type,
		U.username,
		COALESCE(U.profile_photo,'') AS profile_photo,
		U.channel_name,
		(SELECT COUNT(*) FROM views WHERE views.media_id=M.id) AS views_count,
		COALESCE(M.thumbnail,'') AS thumbnail,
		M.created_at
	FROM 
		medias AS M
	JOIN users AS U ON U.id = user_id
	WHERE  ` + mt + `M.title LIKE concat('%',$1::VARCHAR,'%')
	ORDER BY ` + st + `
	LIMIT $2 OFFSET $3;
	`

	rows, err := db.Db.Query(query, search, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var medias []*Media
	for rows.Next() {
		var m Media
		err := rows.Scan(&m.Title, &m.Url, &m.Type, &m.UserName, &m.UserProfile, &m.ChannelName, &m.ViewsCount, &m.Thumbnail, &m.CreatedAt)
		if err != nil {
			return 0, nil, err
		}
		m.CurrentUserLike = NONE
		medias = append(medias, &m)
	}
	return totalPages, medias, nil /// medias will be nil if no media can be find
}

func SearchUserMedias(limit int, offset int, search string, username string, t helper.MediaType, sortType helper.SortType) (int64, []*Media, error) {
	mt := fmt.Sprintf(" M.media_type=%d AND ", t)
	if t == helper.MediaType_ALL {
		mt = ""
	}
	var totalPages int64
	var userId int64
	{
		query := `
		SELECT 
			user_id,COUNT(*)
		FROM 
			medias M
		WHERE ` + mt + `user_id=getUserIdByUsername($1) AND M.title LIKE concat('%',$2::VARCHAR,'%')
		GROUP BY M.user_id;
		`
		rows, err := db.Db.Query(query, search, username)
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
			// return 0, nil, NewModelError("No media found", 404)
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
		M.title,
		M.url,
		M.media_type,
		U.username,
		COALESCE(U.profile_photo,'') AS profile_photo,
		U.channel_name,
		(SELECT COUNT(*) FROM views WHERE views.media_id=M.id) AS views_count,
		COALESCE(M.thumbnail,'') AS thumbnail,
		M.created_at
	FROM 
		medias AS M
	JOIN users AS U ON U.id = user_id
	WHERE ` + mt + `M.user_id=$1 AND M.title LIKE concat('%',$2::VARCHAR,'%')
	ORDER BY ` + st + `
	LIMIT $2 OFFSET $3;
	`

	rows, err := db.Db.Query(query, userId, search, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var medias []*Media
	for rows.Next() {
		var m Media
		err := rows.Scan(&m.Title, &m.Url, &m.Type, &m.UserName, &m.UserProfile, &m.ChannelName, &m.ViewsCount, &m.Thumbnail, &m.CreatedAt)
		if err != nil {
			return 0, nil, err
		}
		m.CurrentUserLike = NONE
		medias = append(medias, &m)
	}
	return totalPages, medias, nil /// medias will be nil if no media can be find
}

func GetMedias(limit int, offset int, t helper.MediaType, sortType helper.SortType) (int64, []*Media, error) {
	mt := fmt.Sprintf(" M.media_type=%d AND ", t)
	if t == helper.MediaType_ALL {
		mt = ""
	}
	var totalPages int64
	{
		query := `
		SELECT 
			COUNT(*)
		FROM 
			medias AS M
		WHERE ` + mt + `1=1;
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
	fmt.Println(st)
	query := `
	SELECT 
		M.title,
		M.url,
		M.media_type,
		U.username,
		COALESCE(U.profile_photo,'') AS profile_photo,
		U.channel_name,
		(SELECT COUNT(*) FROM views WHERE views.media_id=M.id) AS views_count,
		COALESCE(M.thumbnail,'') AS thumbnail,
		M.created_at
	FROM 
		medias AS M
	JOIN users AS U ON U.id = user_id
	WHERE  ` + mt + `1=1
	ORDER BY ` + st + `
	LIMIT $1 OFFSET $2;
	`

	rows, err := db.Db.Query(query, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var medias []*Media
	for rows.Next() {
		var m Media
		err := rows.Scan(&m.Title, &m.Url, &m.Type, &m.UserName, &m.UserProfile, &m.ChannelName, &m.ViewsCount, &m.Thumbnail, &m.CreatedAt)
		if err != nil {
			return 0, nil, err
		}
		m.CurrentUserLike = NONE
		medias = append(medias, &m)
	}
	return totalPages, medias, nil /// medias will be nil if no media can be find
}

func GetUserMedias(limit int, offset int, username string, t helper.MediaType, sortType helper.SortType) (int64, []*Media, error) {
	mt := fmt.Sprintf(" M.media_type=%d AND ", t)
	if t == helper.MediaType_ALL {
		mt = ""
	}
	var totalPages int64
	var userId int64
	{
		query := `
		SELECT 
			user_id,COUNT(*)
		FROM 
			medias M
		WHERE ` + mt + `user_id=getUserIdByUsername($1) 
		GROUP BY M.user_id;
		`
		rows, err := db.Db.Query(query, username)
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
		M.title,
		M.url,
		M.media_type,
		U.username,
		COALESCE(U.profile_photo,'') AS profile_photo,
		U.channel_name,
		(SELECT COUNT(*) FROM views WHERE views.media_id=M.id) AS views_count,
		COALESCE(M.thumbnail,'') AS thumbnail,
		M.created_at
	FROM 
		medias AS M
	JOIN users AS U ON U.id = user_id
	WHERE ` + mt + `M.user_id=$1
	ORDER BY ` + st + `
	LIMIT $2 OFFSET $3;
	`

	rows, err := db.Db.Query(query, userId, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var medias []*Media
	for rows.Next() {
		var m Media
		err := rows.Scan(&m.Title, &m.Url, &m.Type, &m.UserName, &m.UserProfile, &m.ChannelName, &m.ViewsCount, &m.Thumbnail, &m.CreatedAt)
		if err != nil {
			return 0, nil, err
		}
		m.CurrentUserLike = NONE
		medias = append(medias, &m)
	}
	return totalPages, medias, nil /// medias will be nil if no media can be find
}

func GetMediaByUrl(url string) (*Media, error) {
	query := `
	SELECT
		M.id,
		M.title,
		COALESCE(M.text,''),
		M.url,
		M.media_type,
		M.user_id,
		(SELECT COUNT(*) FROM followings AS F WHERE F.following_id = M.user_id) AS subscribers,
		(SELECT COUNT(*) FROM comments WHERE comments.media_id=M.id) AS comments_count,
		(SELECT COUNT(*) FROM likes WHERE likes.media_id=M.id AND likes.is_like=true) AS likes_count,
		(SELECT COUNT(*) FROM likes WHERE likes.media_id=M.id AND likes.is_like=false) AS dislikes_count,
		(SELECT COUNT(*) FROM views WHERE views.media_id=M.id) AS views_count,
		COALESCE((SELECT string_agg(tags.name,',') FROM media_tags JOIN tags ON tags.id=media_tags.tag_id WHERE media_id=M.id),'') AS tags,
		U.username,
		COALESCE(U.profile_photo,'') AS profile_photo,
		U.channel_name,
		COALESCE(M.thumbnail,'') AS thumbnail,
		M.created_at,
		M.updated_at
	FROM 
		medias AS M
	JOIN users AS U 
		ON U.id = user_id
	WHERE 
		m.url=$1;
	`
	rows, err := db.Db.Query(query, url)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, NewModelError("Media with url:`"+url+"` not found", 404)
	}
	var m Media
	err = rows.Scan(
		&m.Id, &m.Title, &m.Text, &m.Url, &m.Type, &m.UserId, &m.UserSubscribesCount,
		&m.CommentsCount, &m.LikesCount, &m.DislikesCount, &m.ViewsCount, &m.Tags,
		&m.UserName, &m.UserProfile, &m.ChannelName, &m.Thumbnail, &m.CreatedAt, &m.UpdatedAt,
	)
	m.CurrentUserLike = NONE

	if err != nil {
		return nil, err
	}
	return &m, nil
}

// auth
func AuthGetMediaByUrl(url string, userId int64) (*Media, error) {
	query := `
	SELECT
		M.id,
		M.title,
		COALESCE(M.text,''),
		M.url,
		M.media_type,
		M.user_id,
		(SELECT COUNT(*) FROM followings AS F WHERE F.following_id = M.user_id) AS subscribers,
		(SELECT COUNT(*) FROM comments WHERE comments.media_id=M.id) AS comments_count,
		(SELECT COUNT(*) FROM likes WHERE likes.media_id=M.id AND likes.is_like=true) AS likes_count,
		(SELECT COUNT(*) FROM likes WHERE likes.media_id=M.id AND likes.is_like=false) AS dislikes_count,
		(SELECT COUNT(*) FROM views WHERE views.media_id=M.id) AS views_count,
		COALESCE((SELECT string_agg(tags.name,',') FROM media_tags JOIN tags ON tags.id=media_tags.tag_id WHERE media_id=M.id),'') AS tags,
		getUserlikeOnMedia($1,M.id) AS is_liked,
		EXISTS(SELECT id FROM followings WHERE follower_id=$1 AND following_id=U.id) AS is_subbed,
		U.username,
		COALESCE(U.profile_photo,'') AS profile_photo,
		U.channel_name,
		COALESCE(M.thumbnail,'') AS thumbnail,
		M.created_at,
		M.updated_at
	FROM 
		medias AS M
	JOIN users AS U 
		ON U.id = user_id
	WHERE 
		m.url=$2;
	`
	rows, err := db.Db.Query(query, userId, url)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, NewModelError("Media with url:`"+url+"` not found", 404)
	}
	var m Media
	err = rows.Scan(
		&m.Id, &m.Title, &m.Text, &m.Url, &m.Type, &m.UserId, &m.UserSubscribesCount, &m.CommentsCount,
		&m.LikesCount, &m.DislikesCount, &m.ViewsCount, &m.Tags, &m.CurrentUserLike, &m.IsCurrentUserSubbedToCreator,
		&m.UserName, &m.UserProfile, &m.ChannelName, &m.Thumbnail, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

/// CreateMedia

func CreateMedia(title string, text string, url string, userId int64, mt helper.MediaType) error {
	query := "INSERT INTO medias (title,text,url,user_id,media_type) VALUES ($1,$2,$3,$4,$5)"
	res, err := db.Db.Exec(query, title, text, url, userId, mt)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			fmt.Printf("%+v\n", *err)
			log.Printf("%+v\n", err.Detail)
			if err.Code == "23505" { // duplicate key
				return &ModelError{Message: err.Detail, Status: 400}
			}
		}
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("failed to create media:" + url)
	}
	return nil
}

//// EditMedia

func EditMedia(url string, title string, text string, userId int64) error {
	query := "UPDATE medias SET title=$1,text=$2,updated_at=$3 WHERE id=(SELECT getMediaIdByUrl($4)) AND user_id=$5;"
	res, err := db.Db.Exec(query, title, text, time.Now().UTC(), url, userId)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%v\n", err.Detail)
			if err.Code == "P0001" {
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
		return NewModelError("Media not belong to user!", 403)
	}
	return nil
}

//// DeleteMedia

func DeleteMedia(url string, userId int64) error {
	query := "DELETE FROM medias WHERE id=(SELECT getMediaIdByUrl($1)) AND user_id=$2;"
	res, err := db.Db.Exec(query, url, userId)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v\n", *err)
			if err.Code == "P0001" {
				return &ModelError{Message: err.Message, Status: 404}
				// return NewModelError("No media with url:`"+url+"` found", 404)
			}
		}
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return NewModelError("Media not belong to user!", 403)
	}
	return nil
}

////////// tags

func AddTagToMedia(mediaUrl string, tagName string, userId int64) error {
	query := `
	INSERT INTO media_tags 
		(media_id,tag_id) 
	SELECT 
		(SELECT getMediaIdByUrl($1)),
		getTagIdByName($2)
	WHERE EXISTS(SELECT user_id FROM medias WHERE user_id=$3 AND id=(SELECT getMediaIdByUrl($4)));
	`
	res, err := db.Db.Exec(query, mediaUrl, tagName, userId, mediaUrl)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v\n", *err)
			if err.Code == "P0001" {
				return &ModelError{Message: err.Message, Status: 404}
			}
			if err.Code == "23505" {
				return NewModelError("tag already added to media!", 400)
			}
		}
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return NewModelError("Media not belong to user!", 403)
	}
	return nil
}

func RemoveTagFromMedia(mediaUrl string, tagName string, userId int64) error {
	query := `
	DELETE FROM media_tags 
	WHERE 
		EXISTS(SELECT id FROM medias WHERE id=getMediaIdByUrl($1) AND user_id=$2)
		AND 
		media_id = (SELECT getMediaIdByUrl($3))
		AND
		tag_id = (SELECT getTagIdByName($4));
	`
	res, err := db.Db.Exec(query, mediaUrl, userId, mediaUrl, tagName)
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
		// return NewModelError("tag is not in media or meida is not created by user!", 400)
		return NewModelError("media or tag not found!", 400)
	}
	return nil
}
