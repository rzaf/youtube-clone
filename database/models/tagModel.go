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
	query := "INSERT INTO tags (name,created_at) VALUES ($1,$2)"
	t := time.Now()
	res, err := db.Db.Exec(query, name, t)
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

//// DeleteLike

// func DeleteTag(tagName int64) error {
// 	query := "DELETE FROM tags WHERE name=$1;"
// 	res, err := db.Db.Exec(query, tagName)
// 	if err != nil {
// 		return err
// 	}
// 	n, err := res.RowsAffected()
// 	if err != nil {
// 		return err
// 	}
// 	if n == 0 {
// 		return NewModelError("tag not found", 404)
// 	}
// 	return nil
// }
