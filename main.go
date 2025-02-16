package main

import (
	"log"

	"github.com/abdulkarimogaji/invoGenius/db"
	"github.com/abdulkarimogaji/invoGenius/server"
)

func main() {
	// connect db
	err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	err = server.StartServer()
	log.Fatal(err)
}
