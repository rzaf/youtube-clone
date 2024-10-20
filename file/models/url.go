package models

import (
	"context"
	"errors"
	"fmt"
	pbHelper "github.com/rzaf/youtube-clone/database/pbs/helper"
	"github.com/rzaf/youtube-clone/file/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// type MediaType uint8

// const (
// 	Video MediaType = iota
// 	Music
// 	Photo
// )

type MediaState uint8

const (
	Pinned MediaState = iota
	Ready
	Removed
)

type MediaOwner uint8

const (
	None MediaOwner = iota
	PhotoMedia
	VideoMedia
	MusicMedia
	ProfilePhoto
	ChannelPhoto
	VideoThumbnail
	PlaylistThumbnail
)

type UrlEntry struct {
	Id        string             `bson:"_id,omitempty"`
	Url       string             `bson:"url"`
	Size      int64              `bson:"size"` /// byte
	Type      pbHelper.MediaType `bson:"media_type"`
	State     MediaState         `bson:"media_state"`
	Owner     MediaOwner         `bson:"media_owner"`
	UserId    int64              `bson:"user_id"`
	CreatedAt time.Time          `bson:"created_at"`
}

func ExistsUrl(url string) bool {
	res := db.UrlsCollection.FindOne(context.Background(), bson.M{
		"url": url,
	})
	return res.Err() != nil
}

func GetUrl(url string) *UrlEntry {
	res := db.UrlsCollection.FindOne(context.Background(), bson.M{
		"url": url,
	})
	if res.Err() == nil {
		var e UrlEntry
		err := res.Decode(&e)
		if err != nil {
			fmt.Println("decoding url failed:", err)
			return nil
		}
		return &e
	}
	return nil
}

func CreateUrl(url string, size int64, t pbHelper.MediaType, userId int64) error {
	_, err := db.UrlsCollection.InsertOne(context.Background(), UrlEntry{
		Url:       url,
		Size:      size,
		Type:      t,
		State:     Pinned,
		Owner:     None,
		UserId:    userId,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		fmt.Println("createUrl failed")
		return err
	}
	return nil
}

func SetUrlOwner(url string, owner MediaOwner) error {
	res := db.UrlsCollection.FindOneAndUpdate(context.Background(), bson.M{"url": url}, bson.M{
		"$set": bson.M{
			"media_owner": owner,
		},
	})
	if res.Err() != nil {
		fmt.Println("SetUrlOwner failed")
		return res.Err()
	}
	return nil
}

func SetUrlState(url string, s MediaState) {
	res := db.UrlsCollection.FindOneAndUpdate(context.Background(), bson.M{"url": url}, bson.M{
		"$set": bson.M{
			"media_state": s,
		},
	})
	if res.Err() != nil {
		fmt.Println("SetUrlState failed")
		panic(res.Err())

	}
}

func GetUserUploadSizeIn24(userId int64) (int64, error) {
	matchStage := bson.D{{"$match", bson.D{
		{"created_at", bson.D{{"$gt", time.Now().UTC().Add(-time.Hour * 24)}}},
		{"user_id", userId},
	}}}
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$user_id"},
			{"total", bson.D{{"$sum", "$size"}}},
		}},
	}
	cursor, err := db.UrlsCollection.Aggregate(context.Background(), mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return 0, err
	}
	if !cursor.Next(context.Background()) {
		return 0, nil
	}
	fmt.Println(cursor.Current)
	res, err := cursor.Current.Values()
	fmt.Println(res)
	// err = cursor.Decode(totalSize)
	if err != nil {
		return 0, err
	}
	totalSize, ok := res[1].AsInt64OK()
	if !ok {
		return 0, errors.New("totalSize is not int64")
	}
	return totalSize, nil
}
