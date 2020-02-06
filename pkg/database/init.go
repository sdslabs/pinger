// Package database contains the connection to database for the status app
// as well as the timeseries db (timescale) for storing the metrics. It
// contains methods and types to interact with the database using an ORM.
package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // PostgreSQL

	"github.com/sdslabs/status/pkg/utils"
)

var (
	dbConf = utils.StatusConf.Database
	db     *gorm.DB
)

// SetupDB sets up the PostgreSQL API.
func SetupDB() error {
	var err error

	connectStr := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s",
		dbConf.Host,
		dbConf.Port,
		dbConf.Username,
		dbConf.Name,
		dbConf.Password)
	if !dbConf.SSLMode {
		connectStr = fmt.Sprintf("%s sslmode=disable", connectStr)
	}

	db, err = gorm.Open("postgres", connectStr)
	if err != nil {
		return err
	}

	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;").Error; err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&User{},
		&Check{},
		&Payload{},
		&Page{},
		&Incident{},
		&Metric{}).Error; err != nil {
		return err
	}

	if err := db.Model(&Payload{}).AddForeignKey(
		"check_id", "checks(id)", "CASCADE", "CASCADE").Error; err != nil {
		return err
	}

	if err := db.Model(&Incident{}).AddForeignKey(
		"page_id", "pages(id)", "CASCADE", "CASCADE").Error; err != nil {
		return err
	}

	if err := db.Exec("CREATE INDEX ON metrics (check_id, start_time DESC);").Error; err != nil {
		return err
	}

	if err := db.Exec(
		"SELECT create_hypertable('metrics', 'start_time', if_not_exists => TRUE, create_default_indexes => FALSE);").Error; err != nil {
		return err
	}

	return nil
}
