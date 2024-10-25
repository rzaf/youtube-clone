package models

import (
	"fmt"
	"github.com/rzaf/youtube-clone/database/db"
	"log"
	"time"

	"github.com/lib/pq"
)

type Following struct {
	Id                int64      `json:"id"`
	FollowerId        int64      `json:"follower_id"`
	FollowingId       int64      `json:"following_id"`
	FollowerUsername  string     `json:"follower_username"`
	FollowingUsername string     `json:"following_username"`
	CreatedAt         *time.Time `json:"created_at"`
}

type FollowerInfo struct {
	UserId      int64
	Username    string
	ChannelName string
	Email       string
	IsVerified  bool
}

//// Get

// func GetFollowersOfUser(followingId int64) []Following {
// 	query := "SELECT id,follower_id FROM followings WHERE following_id=$1;"
// 	rows, err := db.Db.Query(query, followingId)
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	defer rows.Close()
// 	var followers []Following
// 	for !rows.Next() {
// 		var f Following
// 		err = rows.Scan(&f.Id, &f.FollowerId)
// 		if err != nil {
// 			log.Panicln(err.Error())
// 		}
// 	}
// 	return followers
// }

// func GetFollowingsOfUser(followingId int64) []Following {
// 	query := "SELECT id,following_id FROM followings WHERE follower_id=$1;"
// 	rows, err := db.Db.Query(query, followingId)
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	defer rows.Close()
// 	var followers []Following
// 	for !rows.Next() {
// 		var f Following
// 		err = rows.Scan(&f.Id, &f.FollowingId)
// 		if err != nil {
// 			log.Panicln(err.Error())
// 		}
// 	}
// 	return followers
// }

/// Create

func CreateFollowing(followerId int64, followingUsername string) error {
	// query := "INSERT INTO followings (follower_id,following_id,created_at) VALUES ($1,COALESCE((SELECT id FROM users WHERE username=$2),0),$3)"
	query := "INSERT INTO followings (follower_id,following_id) VALUES ($1,getUserIdByUsername($2))"
	res, err := db.Db.Exec(query, followerId, followingUsername)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v\n", *err)
			if err.Code == "23505" { // duplicate ky violation (follower_id,following_id)
				return NewModelError("username:`"+followingUsername+"` is already followd!", 400)
			}
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
		return NewModelError("user not found", 400)
		// return NewModelError("username:`"+followingUsername+"` is already followd!", 400)
	}
	return nil
}

/// Edit

//// Unfollow

func DeleteFollowing(followerId int64, followingUsername string) error {
	query := `
	DELETE FROM 
		followings 
	WHERE
		follower_id =$1 
	AND 
		following_id=getUserIdByUsername($2);
	`
	res, err := db.Db.Exec(query, followerId, followingUsername)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v\n", *err)
			if err.Code == "P0001" {
				return &ModelError{Message: err.Message, Status: 400}
			}
		}
		return err
	}
	n, err := res.RowsAffected()
	fmt.Println(n)
	if err != nil {
		return err
	}
	if n == 0 {
		return NewModelError("username:`"+followingUsername+"` is not followed", 400)
	}
	return nil
}

///// helper functions

func Helper_FollowersById(followingId int64) ([]FollowerInfo, error) {
	query := `
		SELECT
			U.id, 
			U.username, 
			U.email, 
			U.channel_name,
			(U.email_verification IS NULL) is_verified
		FROM 
			followings F 
		JOIN
			users U
		ON
			F.follower_id = U.id
		WHERE 
			F.following_id=$1;
		`
	rows, err := db.Db.Query(query, followingId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var followers []FollowerInfo
	for rows.Next() {
		var f FollowerInfo
		err = rows.Scan(&f.UserId, &f.Username, &f.Email, &f.ChannelName, &f.IsVerified)
		if err != nil {
			return nil, err
		}
		followers = append(followers, f)
	}
	return followers, nil
}
