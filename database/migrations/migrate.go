package migrations

import (
	"fmt"
	"log"
	"youtube-clone/database/db"
)

type Migration interface {
	up() []string
	down() []string
	tableName() string
}

var migrations []Migration = []Migration{
	&Users{},
	&Medias{},
	&Likes{},
	&Comments{},
	&CommentsLikes{},
	&Tags{},
	&MediaTags{},
	&Followings{},
	&Views{},
	&Playlists{},
	&PlaylistMedias{},
	&Histories{},
}

func UpAll() {
	fmt.Println("Starting migrations up")
	for i := 0; i < len(migrations); i++ {
		fmt.Print("migrating " + migrations[i].tableName() + " table .....")
		queries := migrations[i].up()
		for _, query := range queries {
			_, err := db.Db.Exec(query)
			if err != nil {
				log.Fatalln(err)
			}
		}
		fmt.Println("done")
	}
}

func DownAll() {
	fmt.Println("Starting migrations down")
	for i := len(migrations); i >= 0; i-- {
		fmt.Print("migrating " + migrations[i].tableName() + " table .....")
		queries := migrations[i].down()
		for _, query := range queries {
			_, err := db.Db.Exec(query)
			if err != nil {
				log.Fatalln(err)
			}
		}
		fmt.Println("done")
	}
}
