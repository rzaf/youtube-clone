package main

import (
	"log"
	"youtube-clone/database/db"
	"youtube-clone/database/migrations"

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
