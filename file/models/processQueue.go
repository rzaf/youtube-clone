package models

import (
	"context"
	"fmt"
	"time"

	"github.com/rzaf/youtube-clone/file/db"
	"go.mongodb.org/mongo-driver/bson"
)

type QueueEntry struct {
	Id          string     `bson:"_id,omitempty"`
	Name        string     `bson:"name"` // process type
	Url         string     `bson:"url"`
	Retries     int32      `bson:"retries"`
	LastTriedAt *time.Time `bson:"last_tried_at"`
	CreatedAt   time.Time  `bson:"created_at"`
}

type ProcessedEntry struct {
	Id        string    `bson:"_id,omitempty"`
	Name      string    `bson:"name"` // process type
	Url       string    `bson:"url"`
	CreatedAt time.Time `bson:"created_at"`
}

func CreateProcess(name string, url string) error {
	e := QueueEntry{
		Name:        name,
		Url:         url,
		Retries:     0,
		LastTriedAt: nil,
		CreatedAt:   time.Now().UTC(),
	}
	_, err := db.ProcessesQueueCollection.InsertOne(context.Background(), e)
	if err != nil {
		fmt.Println("CreateProccess failed")
		return err
	}
	return nil
}

func SetProcessDone(name string, url string) error {
	_, err := db.ProcessesQueueCollection.DeleteOne(context.Background(), bson.M{
		"name": name,
		"url":  url,
	})
	if err != nil {
		fmt.Println("failed to delete QueueEntry from ProcessesQueue")
		return err
	}

	_, err = db.ProcessedQueueCollection.InsertOne(context.Background(), ProcessedEntry{
		Name:      name,
		Url:       url,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		fmt.Println("failed to insert ProcessedEntry into ProcessedCollection")
		return err
	}
	return nil
}

func SetProcessFailed(url string) error {
	res := db.ProcessesQueueCollection.FindOneAndUpdate(context.Background(), bson.M{"url": url}, bson.M{
		"$set": bson.M{
			"last_tried_at": time.Now().UTC(),
		},
		"$inc": bson.M{
			"retries": 1,
		},
	})
	if res.Err() != nil {
		fmt.Println("SetProcessFailed failed")
		return res.Err()
	}
	return nil
}

func GetRemainingProcesses() ([]QueueEntry, error) {
	var results []QueueEntry
	resCurser, err := db.ProcessesQueueCollection.Find(context.Background(), bson.M{
		"retries": bson.M{"$lte": 3},
	})
	if err != nil {
		fmt.Println("GetRemainingProcesses failed")
		return nil, err
	}
	err = resCurser.All(context.Background(), &results)
	if err != nil {
		fmt.Println("GetRemainingProcesses failed")
		return nil, err
	}
	return results, nil
}
