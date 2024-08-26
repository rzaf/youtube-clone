package helpers

import (
	"fmt"
	"log"
	"os"
	pbHelper "youtube-clone/database/pbs/helper"
)

func FatalIfEmptyVar(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("%s is not set !\n", key)
	}
	return v
}

func SortToString(sortType pbHelper.SortType) (string, error) {
	switch sortType {
	case pbHelper.SortType_MostViewed:
		return " views_count DESC,created_at DESC", nil

	case pbHelper.SortType_LeastViewed:
		return " views_count ASC,created_at DESC", nil

	case pbHelper.SortType_Newest:
		return " created_at DESC ", nil

	case pbHelper.SortType_Oldest:
		return " created_at ASC ", nil

	case pbHelper.SortType_MostSubscribers:
		return " subscribers DESC,created_at DESC", nil

	case pbHelper.SortType_LeastSubscribers:
		return " subscribers ASC,created_at DESC", nil

	case pbHelper.SortType_MostLiked:
		return " likes_count DESC,created_at DESC", nil

	case pbHelper.SortType_LeastLiked:
		return " likes_count ASC,created_at DESC", nil

	case pbHelper.SortType_MostDisiked:
		return " dislikes_count DESC,created_at DESC", nil

	case pbHelper.SortType_LeastDisliked:
		return " dislikes_count ASC,created_at DESC", nil

	case pbHelper.SortType_MostReplied:
		return " replies_count DESC,created_at DESC", nil

	case pbHelper.SortType_LeastReplied:
		return " replies_count ASC,created_at DESC", nil
	}
	return "", fmt.Errorf("type:%d not found", sortType)
}
