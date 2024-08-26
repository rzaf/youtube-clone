package models

// import (
// 	"log"
// 	"time"
// 	"youtube-clone/database/db"
// )

// type MediaTags struct {
// 	Id        int64      `json:"id"`
// 	MediaId   int64      `json:"media_id"`
// 	TagId     int64      `json:"tag_id"`
// 	CreatedAt *time.Time `json:"created_at"`
// }

// //// get

// func GetTagsOfMedia(mediaId int64) []MediaTags {
// 	query := "SELECT id,tag_id FROM tags WHERE media_id=$1;"
// 	rows, err := db.Db.Query(query, mediaId)
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	defer rows.Close()
// 	var tags []MediaTags
// 	for !rows.Next() {
// 		var t MediaTags
// 		err = rows.Scan(&t.Id, &t.TagId)
// 		if err != nil {
// 			log.Panicln(err.Error())
// 		}
// 	}
// 	return tags
// }

// /// CreateTag

// func CreateTagForMedia(tagId int64, mediaId int64) *MediaTags {
// 	query := "INSERT INTO Media_tags (tag_id,user_id,created_at) VALUES ($1,$2,$3)"
// 	t := time.Now()
// 	res, err := db.Db.Exec(query, tagId, mediaId, t)
// 	if err != nil {
// 		// if mySqlErr, ok := err.(*pgz.MySQLError); ok {
// 		// 	if mySqlErr.Number == 1062 {
// 		// 		return nil, errors.New(mySqlErr.Message)
// 		// 	}
// 		// }
// 		log.Panicln(err)
// 	}
// 	id, err := res.LastInsertId()
// 	if err != nil {
// 		log.Panicln(err.Error())
// 	}
// 	return &MediaTags{
// 		Id: id,
// 		// Name:      name,
// 		TagId:     tagId,
// 		MediaId:   mediaId,
// 		CreatedAt: &t,
// 	}
// }

// /// EditTag

// //// DeleteLike

// func DeleteMediaTag(mediaId, tagId int64) {
// 	query := "DELETE FROM media_tags WHERE media_id=$1,tag_id=$2;"
// 	_, err := db.Db.Exec(query, mediaId, tagId)
// 	if err != nil {
// 		log.Panicln(err.Error())
// 	}
// }
