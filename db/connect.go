package db

import (
	"database/sql"

	"github.com/abdulkarimogaji/invoGenius/config"
	_ "github.com/go-sql-driver/mysql"
)

var DB *Queries

func ConnectDB() error {
	connStr := config.C.Database_Uri
	sqlDb, err := sql.Open("mysql", connStr)
	if err != nil {
		return err
	}

	DB = New(sqlDb)
	return nil
}
