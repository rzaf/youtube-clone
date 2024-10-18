package models

import (
	"errors"
	"log"
	"youtube-clone/database/db"

	"github.com/lib/pq"
)

type LikeState uint8

const (
	DISLIKE LikeState = iota
	LIKE
	NONE
)

func (l LikeState) String() string {
	switch l {
	case LIKE:
		return "Liked"
	case DISLIKE:
		return "Disliked"
	case NONE:
		return "None"
	}
	log.Println("ConvertLikeToStr should not return empty !!!")
	return ""
}

////////// MEDIA

/// CreateLike

func CreateMediaLike(userId int64, mediaUrl string, isLike bool) error {
	query := "INSERT INTO likes (user_id,media_id,is_like) VALUES ($1,getMediaIdByUrl($2),$3)"
	res, err := db.Db.Exec(query, userId, mediaUrl, isLike)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v\n", *err)
			if err.Code == "23505" { // duplicate key violation (user_id,media_id)
				return NewModelError("media:`"+mediaUrl+"` is already liked or disliked!", 400)
			}
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
		return errors.New("should not happen")
	}
	return nil
}

//// DeleteLike

func DeleteMediaLike(userId int64, mediaUrl string) error {
	query := "DELETE FROM likes WHERE user_id=$1 AND media_id=getMediaIdByUrl($2);"
	res, err := db.Db.Exec(query, userId, mediaUrl)
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
		return &ModelError{Message: "Media is not liked or disliked", Status: 400}
	}
	return nil
}

////////// COMMENT

/// CreateLike

func CreateCommentLike(userId int64, commentUrl string, isLike bool) error {
	query := "INSERT INTO comments_likes (user_id,comment_id,is_like) VALUES ($1,getCommentIdByUrl($2),$3)"
	res, err := db.Db.Exec(query, userId, commentUrl, isLike)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v\n", *err)
			if err.Code == "23505" { // duplicate key violation (user_id,media_id)
				return NewModelError("comment:`"+commentUrl+"` is already liked or disliked!", 400)
			}
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
		return errors.New("should not happen")
	}
	return nil
}

//// DeleteLike

func DeleteCommentLike(userId int64, commentUrl string) error {
	query := "DELETE FROM comments_likes WHERE user_id=$1 AND comment_id=getCommentIdByUrl($2);"
	res, err := db.Db.Exec(query, userId, commentUrl)
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
		return &ModelError{Message: "Comment is not liked or disliked!", Status: 400}
	}
	return nil
}
