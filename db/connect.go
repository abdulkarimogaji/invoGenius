package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var DB *Queries

func ConnectDB() error {
	connStr := "root:@tcp(localhost:3306)/invoGenius?parseTime=true"
	sqlDb, err := sql.Open("mysql", connStr)
	if err != nil {
		return err
	}

	DB = New(sqlDb)
	return nil
}
