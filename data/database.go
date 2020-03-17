package data

import (
	"log"

	"github.com/asdine/storm"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func InitDB(name string) {
	db, err := storm.Open(name)
	handleError(err, "Unable to initialize database")
	db.Close()
}
