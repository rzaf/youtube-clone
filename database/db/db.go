package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"youtube-clone/database/helpers"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func Connect() {
	host := helpers.FatalIfEmptyVar("DB_HOST")
	port := helpers.FatalIfEmptyVar("DB_PORT")
	user := helpers.FatalIfEmptyVar("DB_USER")
	pass := helpers.FatalIfEmptyVar("DB_PASS")
	name := helpers.FatalIfEmptyVar("DB_NAME")
	ssl := helpers.FatalIfEmptyVar("DB_SSL_MODE")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, name, ssl)
	log.Printf("dsn: %s\n", dsn)
	n := 0
	for n < 10 {
		log.Printf("Trying to connect to postgres database for %d time\n", n+1)
		var err error
		Db, err = sql.Open("postgres", dsn)
		if err == nil {
			log.Printf("pinging db ..... \n")
			pingErr := Db.Ping()
			if pingErr == nil {
				log.Printf("connected \n")
				return
			}
			log.Println("Failed to ping database. pingErr:", pingErr)
		} else {
			log.Println("Failed to connect to database. err:", err)
		}
		n++
		log.Printf("Trying again after 2 second \n")
		time.Sleep(2 * time.Second)
	}
	log.Fatal("connection to database failed after 10 tries !!!")
}

func Disconnect() {
	Db.Close()
}
