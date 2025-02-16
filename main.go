package main

import (
	"log"

	"github.com/abdulkarimogaji/invoGenius/config"
	"github.com/abdulkarimogaji/invoGenius/db"
	"github.com/abdulkarimogaji/invoGenius/server"
)

func main() {
	// load config
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	// connect db
	err = db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	err = server.StartServer()
	log.Fatal(err)
}
