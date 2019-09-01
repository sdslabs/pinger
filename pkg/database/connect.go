package database

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // for sqlite
)

const dbName = "database.sqlite"

func createDBIfNotExist() {
	_, err := os.Stat(dbName)
	if os.IsNotExist(err) {
		db, err := os.Create(dbName)
		if err != nil {
			panic(err)
		}
		defer db.Close()
	} else if err != nil {
		panic(err)
	}
}

// GetSQLDB returns a connection to the sqlite database
func GetSQLDB() (SQLDB, error) {
	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&User{},
		&Check{},
		&Payload{},
		&Page{},
		&Incident{})

	return &sqldb{DB: db}, nil
}

func init() {
	createDBIfNotExist()
}
