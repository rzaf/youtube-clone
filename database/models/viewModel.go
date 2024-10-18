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

func CreateView(ip string, userId int64) *Views {
	query := "INSERT INTO views (ip,user_id) VALUES ($1,$2)"
	res, err := db.Db.Exec(query, ip, userId)
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
		Id:     id,
		Ip:     ip,
		UserId: userId,
	}
}
