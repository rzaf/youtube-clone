package helpers

import (
	"fmt"
	"github.com/rzaf/youtube-clone/notification/pbs/notificationHelperPb"
)

func SeenTypeToSql(seenType notificationHelperPb.SeenType) (string, error) {
	switch seenType {

	case notificationHelperPb.SeenType_Any:
		return "", nil

	case notificationHelperPb.SeenType_Seen:
		return " AND seen_at IS NOT NULL ", nil

	case notificationHelperPb.SeenType_NotSeen:
		return " AND seen_at IS NULL ", nil

	}
	return "", fmt.Errorf("type:%d not found", seenType)
}

func SortToSql(sortType notificationHelperPb.SortType) (string, error) {
	switch sortType {

	case notificationHelperPb.SortType_Newest:
		return " created_at DESC ", nil

	case notificationHelperPb.SortType_Oldest:
		return " created_at ASC ", nil

	}
	return "", fmt.Errorf("type:%d not found", sortType)
}
