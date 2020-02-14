// Package database contains the connection to database for the status app
// as well as the timeseries db (timescale) for storing the metrics. It
// contains methods and types to interact with the database using an ORM.
package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // PostgreSQL

	"github.com/sdslabs/status/pkg/config"
	"github.com/sdslabs/status/pkg/defaults"
	"github.com/sdslabs/status/pkg/metrics"
)

var db *gorm.DB

func initFromProvider(conf *metrics.ProviderConfig) error {
	var err error

	connectStr := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s",
		conf.Host,
		conf.Port,
		conf.Username,
		conf.DBName,
		conf.Password)
	if !conf.SSLMode {
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

func createStandaloneUser() (uint, error) {
	// This creates a user corresponding to which all the standalone metrics will be saved.
	// This user need not be a valid email address. Later on the metrics can be exported
	// corresponding to the checks this user creates.
	u, err := CreateUser(&User{
		Email: defaults.StandaloneUserEmail,
		Name:  defaults.StandaloneUserName,
	})
	if err != nil {
		return 0, err
	}
	return u.ID, nil
}

// SetupDB sets up the PostgreSQL API.
func SetupDB(conf *config.StatusConfig) error {
	dbConf := &metrics.ProviderConfig{
		Backend:  metrics.TimeScaleProviderType,
		Host:     conf.Application.Database.Host,
		Port:     conf.Application.Database.Port,
		Username: conf.Application.Database.Username,
		Password: conf.Application.Database.Password,
		DBName:   conf.Application.Database.Name,
		SSLMode:  conf.Application.Database.SSLMode,
	}
	if err := initFromProvider(dbConf); err != nil {
		return err
	}
	if _, err := createStandaloneUser(); err != nil {
		return err
	}
	return nil
}
