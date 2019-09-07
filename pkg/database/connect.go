package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // PostgreSQL
	"github.com/sdslabs/status/pkg/utils"
)

var (
	dbConf = utils.StatusConf.Database
	// DBConn for sending API Queries
	DBConn SQLDB
)

// GetSQLDB returns a connection to the sqlite database
func GetSQLDB() (SQLDB, error) {
	connectStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s", dbConf.Host, dbConf.Port, dbConf.Username, dbConf.Name, dbConf.Password)
	if !dbConf.SSLMode {
		connectStr = fmt.Sprintf("%s sslmode=disable", connectStr)
	}
	db, err := gorm.Open("postgres", connectStr)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&User{},
		&Check{},
		&Payload{},
		&Page{},
		&Incident{})

	db.Model(&Payload{}).AddForeignKey("check_id", "checks(id)", "CASCADE", "CASCADE")
	db.Model(&Incident{}).AddForeignKey("page_id", "pages(id)", "CASCADE", "CASCADE")

	return &sqldb{DB: db}, nil
}

func init() {
	var err error
	DBConn, err = GetSQLDB()
	if err != nil {
		panic(err)
	}
}
