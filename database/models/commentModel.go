package models

import (
	"fmt"
	"log"
	"math"
	"time"
	"youtube-clone/database/db"
	"youtube-clone/database/helpers"
	"youtube-clone/database/pbs/helper"

	"github.com/lib/pq"
)

type Comment struct {
	// Id        int64      `json:"-"`
	Url       string     `json:"url"`
	Text      string     `json:"content"`
	UserId    int64      `json:"-"`
	MediaId   int64      `json:"-"`
	CommentId int64      `json:"-"`
	ReplyUrl  string     `json:"-"`
	ReplyText string     `json:"-"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	/// extra columns from other tables
	RepliesCount     int64  `json:"replies_count"`
	LikesCount       int64  `json:"likes_count"`    ///todo
	DislikesCount    int64  `json:"dislikes_count"` ///todo
	ReplyUsername    string `json:"reply_username"`
	ReplyUserProfile string `json:"reply_user_profile"`
	UserName         string `json:"username"`
	UserProfile      string `json:"userprofile"`
	MediaUrl         string `json:"media_url"`
	CurrentUserLike  LikeState
	///
	MediaType           helper.MediaType `json:"media_type"` //0:video 1:music 2:photo
	MediaCreator        string           `json:"media_creator"`
	MediaCreatorProfile string           `json:"media_creator_profile"`
	MediaTitle          string           `json:"media_title"`
	MediaThumbnail      string           `json:"media_thumbnail"`
	MediaCreatedAt      *time.Time       `json:"media_createdAt"`
}

//// GetComment

func AuthGetComment(commentUrl string, userId int64) (*Comment, error) {
	query := `
	SELECT
		C2.text AS reply_text,
		COALESCE(C2.url,'') AS reply_url,
		COALESCE(U3.username,'') AS reply_username,
		COALESCE(U3.profile_photo,'') AS reply_user_profile,
		C.url,
		C.text,
		M.url,
		M.title,
		M.media_type,
		COALESCE(M.thumbnail,'') AS media_thumbnail,
		M.created_at,
		U2.username,
		COALESCE(U2.profile_photo,'') AS profile_photo2,
		U.username,
		COALESCE(U.profile_photo,'') AS profile_photo1,
		(SELECT COUNT(*) FROM comments WHERE comment_id=c.id) AS replies_count,
		(SELECT COUNT(*) FROM comments_likes WHERE comments_likes.comment_id=C.id AND comments_likes.is_like=true) AS likes_count,
		(SELECT COUNT(*) FROM comments_likes WHERE comments_likes.comment_id=C.id AND comments_likes.is_like=false) AS dislikes_count,
		getUserlikeOnComment(C.id,$1) AS is_liked,
		C.created_at,
		C.updated_at
	FROM
		comments C
	LEFT JOIN comments C2
		ON C.comment_id=C2.id
	JOIN users U
		ON U.id = C.user_id
	JOIN medias M
		ON M.id=C.media_id
	JOIN users U2
		ON U2.id=M.user_id
	LEFT JOIN users U3
		ON U3.id=C2.user_id
	WHERE
		C.url=$2 ;
	`
	rows, err := db.Db.Query(query, userId, commentUrl)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, NewModelError("comment with url:`"+commentUrl+"`not found", 404)
	}
	var c Comment
	err = rows.Scan(&c.ReplyText, &c.ReplyUrl, &c.ReplyUsername, &c.ReplyUserProfile, &c.Url, &c.Text, &c.MediaUrl, &c.MediaTitle, &c.MediaType, &c.MediaThumbnail, &c.MediaCreatedAt, &c.MediaCreator, &c.MediaCreatorProfile, &c.UserName, &c.UserProfile, &c.RepliesCount, &c.LikesCount, &c.DislikesCount, &c.CurrentUserLike, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func GetComment(commentUrl string) (*Comment, error) {
	fmt.Println(commentUrl)
	query := `
	SELECT
		C2.text AS reply_text,
		COALESCE(C2.url,'') AS reply_url,
		COALESCE(U3.username,'') AS reply_username,
		COALESCE(U3.profile_photo,'') AS reply_user_profile,
		C.url,
		C.text,
		M.url,
		M.title,
		M.media_type,
		COALESCE(M.thumbnail,'') AS media_thumbnail,
		M.created_at,
		U2.username,
		COALESCE(U2.profile_photo,'') AS profile_photo2,
		U.username,
		COALESCE(U.profile_photo,'') AS profile_photo1,
		(SELECT COUNT(*) FROM comments WHERE comment_id=c.id) AS replies_count,
		(SELECT COUNT(*) FROM comments_likes WHERE comments_likes.comment_id=C.id AND comments_likes.is_like=true) AS likes_count,
		(SELECT COUNT(*) FROM comments_likes WHERE comments_likes.comment_id=C.id AND comments_likes.is_like=false) AS dislikes_count,
		C.created_at,
		C.updated_at
	FROM
		comments C
	LEFT JOIN comments C2
		ON C.comment_id=C2.id
	JOIN users U
		ON U.id = C.user_id
	JOIN medias M
		ON M.id=C.media_id
	JOIN users U2
		ON U2.id=M.user_id
	LEFT JOIN users U3
		ON U3.id=C2.user_id
	WHERE
		C.url=$1 ;
	`
	rows, err := db.Db.Query(query, commentUrl)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, NewModelError("comment with url:`"+commentUrl+"`not found", 404)
	}
	var c Comment
	err = rows.Scan(&c.ReplyText, &c.ReplyUrl, &c.ReplyUsername, &c.ReplyUserProfile, &c.Url, &c.Text, &c.MediaUrl, &c.MediaTitle, &c.MediaType, &c.MediaThumbnail, &c.MediaCreatedAt, &c.MediaCreator, &c.MediaCreatorProfile, &c.UserName, &c.UserProfile, &c.RepliesCount, &c.LikesCount, &c.DislikesCount, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	c.CurrentUserLike = NONE
	return &c, nil
}

// auth
func AuthGetCommentsOfMedia(mediaUrl string, limit int, offset int, userId int64, sortType helper.SortType) (int64, []*Comment, error) {
	st, err := helpers.SortToString(sortType)
	if err != nil {
		return 0, nil, err
	}
	var totalPages int64
	var mediaId int64
	{
		query := `
		SELECT 
			media_id,COUNT(*)
		FROM 
			comments
		WHERE media_id=getMediaIdByUrl($1) AND comment_id IS NULL
		GROUP BY media_id
		`
		rows, err := db.Db.Query(query, mediaUrl)
		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				log.Printf("%+v\n", *err)
				if err.Code == "P0001" {
					return 0, nil, &ModelError{Message: err.Message, Status: 400}
				}
			}
			return 0, nil, err
		}
		defer rows.Close()
		if !rows.Next() {
			return 0, nil, nil
		}
		err = rows.Scan(&mediaId, &totalPages)
		if err != nil {
			return 0, nil, err
		}
		totalPages = int64(math.Ceil(float64(totalPages) / float64(limit)))
	}

	var comments []*Comment
	query := `
		SELECT 
			C.Url,
			C.text,
			U.username,
			COALESCE(U.profile_photo,'') AS profile_photo,
			(SELECT COUNT(*) FROM comments WHERE comment_id=c.id) AS replies_count,
			(SELECT COUNT(*) FROM comments_likes WHERE comments_likes.comment_id=C.id AND comments_likes.is_like=true) AS likes_count,
			(SELECT COUNT(*) FROM comments_likes WHERE comments_likes.comment_id=C.id AND comments_likes.is_like=false) AS dislikes_count,
			getUserlikeOnComment(C.id,$1) AS is_liked,
			C.created_at,
			C.updated_at
		FROM 
			comments C
		JOIN users U 
		ON U.id = C.user_id
		WHERE 
			media_id=$2 
			AND 
			C.comment_id IS NULL 
		ORDER BY ` + st + `
		LIMIT $3 OFFSET $4;
		`
	rows, err := db.Db.Query(query, userId, mediaId, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var c Comment
		err = rows.Scan(&c.Url, &c.Text, &c.UserName, &c.UserProfile, &c.RepliesCount, &c.LikesCount, &c.DislikesCount, &c.CurrentUserLike, &c.CreatedAt, &c.UpdatedAt)
		fmt.Println(err)
		if err != nil {
			return 0, nil, err
		}
		comments = append(comments, &c)
	}
	return totalPages, comments, nil
}

func GetCommentsOfMedia(mediaUrl string, limit int, offset int, sortType helper.SortType) (int64, []*Comment, error) {
	st, err := helpers.SortToString(sortType)
	if err != nil {
		return 0, nil, err
	}
	var totalPages int64
	var mediaId int64
	{
		query := `
		SELECT 
			media_id,COUNT(*)
		FROM 
			comments
		WHERE media_id=getMediaIdByUrl($1) AND comment_id IS NULL
		GROUP BY media_id
		`
		rows, err := db.Db.Query(query, mediaUrl)
		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				log.Printf("%+v\n", *err)
				if err.Code == "P0001" {
					return 0, nil, &ModelError{Message: err.Message, Status: 400}
				}
			}
			return 0, nil, err
		}
		defer rows.Close()
		if !rows.Next() {
			return 0, nil, nil
		}
		err = rows.Scan(&mediaId, &totalPages)
		if err != nil {
			return 0, nil, err
		}
		totalPages = int64(math.Ceil(float64(totalPages) / float64(limit)))
	}

	var comments []*Comment
	query := `
		SELECT 
			C.Url,
			C.text,
			U.username,
			COALESCE(U.profile_photo,'') AS profile_photo,
			(SELECT COUNT(*) FROM comments WHERE comment_id=c.id) AS replies_count,
			(SELECT COUNT(*) FROM comments_likes WHERE comments_likes.comment_id=C.id AND comments_likes.is_like=true) AS likes_count,
			(SELECT COUNT(*) FROM comments_likes WHERE comments_likes.comment_id=C.id AND comments_likes.is_like=false) AS dislikes_count,
			C.created_at,
			C.updated_at
		FROM 
			comments C
		JOIN users U 
		ON U.id = C.user_id
		WHERE 
			media_id=$1 
			AND 
			C.comment_id IS NULL 
		ORDER BY ` + st + `
		LIMIT $2 OFFSET $3;
		`
	rows, err := db.Db.Query(query, mediaId, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var c Comment
		err = rows.Scan(&c.Url, &c.Text, &c.UserName, &c.UserProfile, &c.RepliesCount, &c.LikesCount, &c.DislikesCount, &c.CreatedAt, &c.UpdatedAt)
		c.CurrentUserLike = NONE
		fmt.Println(err)
		if err != nil {
			return 0, nil, err
		}
		comments = append(comments, &c)
	}
	return totalPages, comments, nil
}

func GetRepliesOfComment(commentUrl string, limit int, offset int, sortType helper.SortType) (int64, []*Comment, error) {
	st, err := helpers.SortToString(sortType)
	if err != nil {
		return 0, nil, err
	}
	var totalPages int64
	var commentId int64
	{
		query := `
		SELECT 
			comment_id,COUNT(*)
		FROM 
			comments
		WHERE comment_id=getCommentIdByUrl($1) 
		GROUP BY comment_id;
		`
		rows, err := db.Db.Query(query, commentUrl)
		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				log.Printf("%+v\n", *err)
				if err.Code == "P0001" {
					return 0, nil, &ModelError{Message: err.Message, Status: 400}
				}
			}
			return 0, nil, err
		}
		defer rows.Close()
		if !rows.Next() {
			return 0, nil, nil
		}
		err = rows.Scan(&commentId, &totalPages)
		if err != nil {
			return 0, nil, err
		}
		totalPages = int64(math.Ceil(float64(totalPages) / float64(limit)))
	}

	// query := "SELECT id,text,user_id,media_id FROM comments WHERE media_id=$1 AND user_id=$2 AND comment_id=$3;"
	query := `
	SELECT
		C.Url,
		C.text,
        U.username,
		COALESCE(U.profile_photo,'') AS profile_photo,
        (SELECT COUNT(*) FROM comments WHERE comment_id=c.id) AS replies_count,
		(SELECT COUNT(*) FROM comments_likes WHERE comments_likes.comment_id=C.id AND comments_likes.is_like=true) AS likes_count,
		(SELECT COUNT(*) FROM comments_likes WHERE comments_likes.comment_id=C.id AND comments_likes.is_like=false) AS dislikes_count,
		C.created_at,
		C.updated_at
	FROM 
		comments C
    JOIN users U 
    ON U.id = C.user_id
	WHERE  
        C.comment_id = $1 
	ORDER BY ` + st + `
	LIMIT $2 OFFSET $3
	`
	rows, err := db.Db.Query(query, commentId, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var comments []*Comment
	for rows.Next() {
		var c Comment
		err = rows.Scan(&c.Url, &c.Text, &c.UserName, &c.UserProfile, &c.RepliesCount, &c.LikesCount, &c.DislikesCount, &c.CreatedAt, &c.UpdatedAt)
		c.CurrentUserLike = NONE
		if err != nil {
			return 0, nil, err
		}
		comments = append(comments, &c)
	}
	return totalPages, comments, nil
}

func AuthGetRepliesOfComment(commentUrl string, limit int, offset int, userId int64, sortType helper.SortType) (int64, []*Comment, error) {
	st, err := helpers.SortToString(sortType)
	if err != nil {
		return 0, nil, err
	}
	var totalPages int64
	var commentId int64
	{
		query := `
		SELECT 
			comment_id,COUNT(*)
		FROM 
			comments
		WHERE comment_id=getCommentIdByUrl($1) 
		GROUP BY comment_id
		`
		rows, err := db.Db.Query(query, commentUrl)
		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				log.Printf("%+v\n", *err)
				if err.Code == "P0001" {
					return 0, nil, &ModelError{Message: err.Message, Status: 400}
				}
			}
			return 0, nil, err
		}
		defer rows.Close()
		if !rows.Next() {
			return 0, nil, nil
		}
		err = rows.Scan(&commentId, &totalPages)
		if err != nil {
			return 0, nil, err
		}
		totalPages = int64(math.Ceil(float64(totalPages) / float64(limit)))
	}

	// query := "SELECT id,text,user_id,media_id FROM comments WHERE media_id=$1 AND user_id=$2 AND comment_id=$3;"
	query := `
	SELECT
		C.Url,
		C.text,
        U.username,
		COALESCE(U.profile_photo,'') AS profile_photo,
        (SELECT COUNT(*) FROM comments WHERE comment_id=c.id) AS replies_count,
		(SELECT COUNT(*) FROM comments_likes WHERE comments_likes.comment_id=C.id AND comments_likes.is_like=true) AS likes_count,
		(SELECT COUNT(*) FROM comments_likes WHERE comments_likes.comment_id=C.id AND comments_likes.is_like=false) AS dislikes_count,
		getUserlikeOnComment(C.id,$1) AS is_liked,
		C.created_at,
		C.updated_at
	FROM 
		comments C
    JOIN users U 
    ON U.id = C.user_id
	WHERE  
        C.comment_id = $2 
	ORDER BY ` + st + `
	LIMIT $3 OFFSET $4
	`
	rows, err := db.Db.Query(query, userId, commentId, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var comments []*Comment
	for rows.Next() {
		var c Comment
		err = rows.Scan(&c.Url, &c.Text, &c.UserName, &c.UserProfile, &c.RepliesCount, &c.LikesCount, &c.DislikesCount, &c.CurrentUserLike, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return 0, nil, err
		}
		comments = append(comments, &c)
	}
	return totalPages, comments, nil
}

// auth
func GetCommentsOfUserInMedia(userId int64, mediaUrl string, limit int, offset int, sortType helper.SortType) (int64, []*Comment, error) {
	st, err := helpers.SortToString(sortType)
	if err != nil {
		return 0, nil, err
	}
	var totalPages int64
	var mediaId int64
	{
		query := `
		SELECT 
			media_id,COUNT(*)
		FROM 
			comments
		WHERE media_id=getMediaIdByUrl($1) AND user_id=$2 AND comment_id IS NULL
		GROUP BY media_id;
		`
		rows, err := db.Db.Query(query, mediaUrl, userId)
		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				log.Printf("%+v\n", *err)
				if err.Code == "P0001" {
					return 0, nil, &ModelError{Message: err.Message, Status: 400}
				}
			}
			return 0, nil, err
		}
		defer rows.Close()
		if !rows.Next() {
			return 0, nil, nil
		}
		err = rows.Scan(&mediaId, &totalPages)
		if err != nil {
			return 0, nil, err
		}
		totalPages = int64(math.Ceil(float64(totalPages) / float64(limit)))
	}

	query := `
	SELECT
		C.Url,
		C.text,
        U.username,
		COALESCE(U.profile_photo,'') AS profile_photo,
        (SELECT COUNT(*) FROM comments WHERE comment_id=c.id) AS replies_count,
		(SELECT COUNT(*) FROM comments_likes WHERE comments_likes.comment_id=C.id AND comments_likes.is_like=true) AS likes_count,
		(SELECT COUNT(*) FROM comments_likes WHERE comments_likes.comment_id=C.id AND comments_likes.is_like=false) AS dislikes_count,
		C.created_at,
		C.updated_at
	FROM
		comments C
    JOIN users U
    ON U.id = C.user_id
	WHERE
		user_id=$1
		AND
		media_id=(SELECT id FROM medias WHERE url=$2)
        AND
        C.comment_id IS NULL 
	ORDER BY ` + st + `
	LIMIT $3 OFFSET $4;
	`
	rows, err := db.Db.Query(query, userId, mediaUrl, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var comments []*Comment
	for rows.Next() {
		var c Comment
		err = rows.Scan(&c.Url, &c.Text, &c.UserName, &c.UserProfile, &c.RepliesCount, &c.LikesCount, &c.DislikesCount, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return 0, nil, err
		}
		comments = append(comments, &c)
	}
	return totalPages, comments, nil
}

// auth
func GetAllCommentsOfUser(userId int64, limit int, offset int, sortType helper.SortType) (int64, []*Comment, error) {
	st, err := helpers.SortToString(sortType)
	if err != nil {
		return 0, nil, err
	}
	var totalPages int64
	{
		query := `
		SELECT 
			COUNT(*)
		FROM 
			comments
		WHERE user_id=$1;
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

	// query := "SELECT id,text,user_id,media_id FROM comments WHERE media_id=$1 AND user_id=$2 AND comment_id=NULL;"
	query := `
	SELECT 
		C.Url AS comment_url,
		C.text,
		M.url  AS media_url,
		--C.comment_id,
		(
			SELECT users.username 
			FROM comments 
			JOIN users ON users.id=comments.user_id 
			WHERE comments.id=C.comment_id 
		) AS reply_username,
		(SELECT COUNT(*) FROM comments WHERE comment_id=C.id) replies_count
		getUserlikeOnComment(C.id,$1) AS is_liked,
	FROM 
		comments C
	JOIN medias M ON M.id = C.media_id
	WHERE  
		C.user_id = $2
	ORDER BY ` + st + `
	LIMIT $3 OFFSET $4;
	`
	rows, err := db.Db.Query(query, userId, userId, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var comments []*Comment
	for rows.Next() {
		var c Comment
		err = rows.Scan(&c.Url, &c.Text, &c.MediaUrl, &c.ReplyUsername, &c.RepliesCount, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			log.Panicln(err.Error())
		}
		comments = append(comments, &c)
	}
	return totalPages, comments, nil
}

// / CreateComment
func CreateComment(text string, userId int64, mediaUrl string) (*Comment, error) {
	query := `
	INSERT INTO 
		comments (
			text,
			media_id,
			comment_id,
			user_id,
			url
		) VALUES (
		$1,
		(SELECT id FROM medias WHERE url=$2),
		NULL,
		$3,
		$4);
	`
	url := generateSecureToken(16)
	res, err := db.Db.Exec(query, text, mediaUrl, userId, url)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v\n", *err)
			switch err.Column {
			case "media_id":
				return nil, &ModelError{Message: "Incorrect media url:`" + mediaUrl + "`", Status: 400}
				// case "comment_id":
			}
			// if err.Code == "23505" {
			// 	return nil, &ModelError{Message: err.Detail, Status: 400}
			// }
		}
		return nil, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return nil, err
	}
	return &Comment{
		Url: url,
	}, nil
}

// TODO: should check reply has the same media_id as comment
func CreateReply(text string, userId int64, mediaUrl string, replyUrl string) (*Comment, error) {
	query := `
	INSERT INTO 
		comments (
			text,
			media_id,
			comment_id,
			user_id,
			url
		) VALUES (
		$1,
		(SELECT id FROM medias WHERE url=$2),
		COALESCE((SELECT id FROM comments WHERE url=$3 AND media_id=(SELECT id FROM medias WHERE url=$4)),0),
		$5,
		$6);
	`
	url := generateSecureToken(16)
	res, err := db.Db.Exec(query, text, mediaUrl, replyUrl, mediaUrl, userId, url)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v\n", *err)
			switch err.Column {
			case "media_id":
				return nil, &ModelError{Message: "Incorrect media url:`" + mediaUrl + "`", Status: 400}
				// case "comment_id":
			}
			if err.Code == "23503" {
				return nil, &ModelError{Message: "Incorrect reply url:`" + replyUrl + "`", Status: 400}
				// return nil, &ModelError{Message: err.Detail, Status: 400}
			}
		}
		return nil, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return nil, err
	}
	return &Comment{
		Url: url,
	}, nil
}

/// EditComment

func EditComment(url string, text string, userId int64) error {
	query := "UPDATE comments SET text=$1,updated_at=$2 WHERE url=$3 AND user_id=$4;"
	var err error
	r, err := db.Db.Exec(query, text, time.Now().UTC(), url, userId)
	if err != nil {
		// if err, ok := err.(*pq.Error); ok {
		// 	log.Printf("%v\n", err.Detail)
		// 	if err.Code == "23505" {
		// 		return &ModelError{Message: err.Detail, Status: 400}
		// 	}
		// }
		return err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return NewModelError("comment not found or not belong to user", 400)
	}
	return nil
}

//// DeleteLike

func DeleteComment(url string, userId int64) error {
	query := "DELETE FROM comments WHERE url=$1 AND user_id=$2;"
	r, err := db.Db.Exec(query, url, userId)
	if err != nil {
		return err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return NewModelError("comment not found or not belong to user", 400)
	}
	return nil
}
