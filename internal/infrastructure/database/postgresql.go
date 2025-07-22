package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type ConnectionDB struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	// Sslmode  string
}

func ConnectDB(info ConnectionDB) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", info.Host, info.Port, info.User, info.Password, info.Name)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
