package main

import (
	"github.com/rzaf/youtube-clone/database/db"
	"github.com/rzaf/youtube-clone/database/migrations"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load(".env")
}
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	db.Connect()
	defer db.Disconnect()
	migrations.UpAll()
}
