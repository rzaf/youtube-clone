package models

import (
	"github.com/google/uuid"
	"github.com/rzaf/youtube-clone/notification/db"
	"github.com/rzaf/youtube-clone/notification/helpers"
	"github.com/rzaf/youtube-clone/notification/pbs/notificationHelperPb"
	"log"
	"math"
	"time"

	"github.com/lib/pq"
)

type Notification struct {
	Id        string     `json:"Id"`
	UserId    int64      `json:"-"`
	Title     string     `json:"Title"`
	Message   string     `json:"Message"`
	Data      string     `json:"Data"`
	SeenAt    *time.Time `json:"SeenAt"`
	CreatedAt *time.Time `json:"CreatedAt"`
	UpdatedAt *time.Time `json:"UpdatedAt"`
}

func GetNotification(notificationId string, userId int64) (*Notification, error) {
	query := `
	SELECT
		id,
		user_id,
		title,
		message,
		seen_at,
		created_at,
		updated_at
	FROM 
		notifications
	WHERE
		id::text=$1 AND user_id=$2 ;
	`
	rows, err := db.Db.Query(query, notificationId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, NewModelError("notification with id:`"+notificationId+"`not found", 404)
	}
	var n Notification
	err = rows.Scan(&n.Id, &n.UserId, &n.Title, &n.Message, &n.SeenAt, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func GetAllNotification(userId int64, limit int, offset int, sortType notificationHelperPb.SortType, seenType notificationHelperPb.SeenType) (int64, []*Notification, error) {
	st, err := helpers.SortToSql(sortType)
	if err != nil {
		return 0, nil, err
	}

	seenAtCondition, err := helpers.SeenTypeToSql(seenType)
	if err != nil {
		return 0, nil, err
	}

	var totalPages int64
	{
		query := `
		SELECT 
			COUNT(*)
		FROM 
			notifications
		WHERE user_id=$1 ` + seenAtCondition + `;
		`
		rows, err := db.Db.Query(query, userId)
		if err != nil {
			return 0, nil, err
		}
		defer rows.Close()
		if !rows.Next() {
			return 0, nil, nil
		}
		err = rows.Scan(&totalPages)
		if err != nil {
			return 0, nil, err
		}
		totalPages = int64(math.Ceil(float64(totalPages) / float64(limit)))
	}

	query := `
	SELECT 
		id,
		user_id,
		title,
		message,
		seen_at,
		created_at,
		updated_at
	FROM 
		notifications C
	WHERE  
		C.user_id = $1 ` + seenAtCondition + `
	ORDER BY ` + st + `
	LIMIT $2 OFFSET $3;
	`
	rows, err := db.Db.Query(query, userId, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var notifications []*Notification
	for rows.Next() {
		var n Notification
		err = rows.Scan(&n.Id, &n.UserId, &n.Title, &n.Message, &n.SeenAt, &n.CreatedAt, &n.UpdatedAt)
		if err != nil {
			log.Panicln(err.Error())
		}
		notifications = append(notifications, &n)
	}
	return totalPages, notifications, nil
}

func CreateNotification(userId int64, title string, message string) (*Notification, error) {
	query := `
	INSERT INTO 
		notifications (
			id,
			user_id,
			title,
			message
		) VALUES (
			$1,$2,$3,$4
		)
	;
	`
	notificationId := uuid.NewString()
	res, err := db.Db.Exec(query, notificationId, userId, title, message)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v\n", *err)
		}
		return nil, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return nil, err
	}
	return &Notification{
		Id: notificationId,
	}, nil
}

func ReadAllNotification(userId int64) error {
	query := "UPDATE notifications SET seen_at=$1,updated_at=$1 WHERE user_id=$2;"
	var err error
	r, err := db.Db.Exec(query, time.Now().UTC(), userId)
	if err != nil {
		return err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return NewModelError("no notification found", 400)
	}
	return nil
}

func ReadNotification(id string, userId int64) error {
	query := "UPDATE notifications SET seen_at=$1,updated_at=$1 WHERE id=$2 AND user_id=$3;"
	var err error
	r, err := db.Db.Exec(query, time.Now().UTC(), id, userId)
	if err != nil {
		return err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return NewModelError("no notification found", 400)
	}
	return nil
}

func DeleteNotification(id string, userId int64) error {
	query := "DELETE FROM notifications WHERE id=$1 AND user_id=$2;"
	r, err := db.Db.Exec(query, id, userId)
	if err != nil {
		return err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return NewModelError("notification not found or not belong to user", 400)
	}
	return nil
}
