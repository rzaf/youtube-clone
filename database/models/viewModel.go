package models

import (
	"log"
	"time"
	"youtube-clone/database/db"
)

type Views struct {
	Id        int64      `json:"id"`
	Ip        string     `json:"ip"`
	MediaId   int64      `json:"media_id"`
	UserId    int64      `json:"user_id"`
	CreatedAt *time.Time `json:"created_at"`
}

//// GetComment

// func GetViewCountOfMedia(mediaId int64) int64 {
// 	query := "SELECT COUNT(*) FROM views WHERE media_id=$1;"
// 	res, err := db.Db.Exec(query, mediaId)
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	r, err := res.LastInsertId()
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	return r
// }

/// CreateComment

func CreateView(ip string, userId int64) *Views {
	query := "INSERT INTO views (ip,user_id,created_at) VALUES ($1,$2)"
	t := time.Now()
	res, err := db.Db.Exec(query, ip, userId, t)
	if err != nil {
		// if mySqlErr, ok := err.(*pgz.MySQLError); ok {
		// 	if mySqlErr.Number == 1062 {
		// 		return nil, errors.New(mySqlErr.Message)
		// 	}
		// }
		log.Panicln(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Panicln(err.Error())
	}
	return &Views{
		Id:        id,
		Ip:        ip,
		UserId:    userId,
		CreatedAt: &t,
	}
}

/// EditComment

//// Unfollow
