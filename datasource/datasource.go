package datasource

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func InitConnection() *sqlx.DB {
	log.Println("CONNECTING TO DATABASE")
	db, err := sqlx.Connect("mysql", "root:root1234@(localhost:3306)/pemilu?parseTime=true")
	if err != nil {
		log.Fatalln("Failed to connect to database,    ", err)
	}
	return db
}
