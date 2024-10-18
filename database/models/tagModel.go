package models

import (
	"errors"
	"time"
	"youtube-clone/database/db"
)

type Tag struct {
	Id        int64      `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	//extra
	// MediaCount int64
}

// func GetTagId(name string) (int64, error) {
// 	query := "SELECT id FROM TAGS WHERE name = $1;"
// 	rows, err := db.Db.Query(query, name)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer rows.Close()
// 	if !rows.Next() {
// 		return -1, nil
// 	}
// 	var id int64
// 	err = rows.Scan(&id)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return id, nil
// }

func CreateTag(name string) error {
	query := "INSERT INTO tags (name) VALUES ($1);"
	res, err := db.Db.Exec(query, name)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("failed to create playlist")
	}
	return nil
}
