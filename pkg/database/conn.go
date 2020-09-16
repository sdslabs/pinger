// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config can be used to create a connection with the database.
type Config interface {
	GetName() string // Name of the database

	GetHost() string // Host of the database
	GetPort() uint16 // Port of the database

	GetUsername() string // Username of the database
	GetPassword() string // Password of the database

	IsSSLMode() bool // Should connect using SSL
}

// Conn is the database connection which can be used to access the API to
// interact with the database.
type Conn struct{ db *gorm.DB }

// NewConn creates a new connection with the database.
func NewConn(conf Config) (*Conn, error) {
	connStr := fmt.Sprintf(
		`host=%s port=%d user=%s dbname=%s password=%s`,
		conf.GetHost(),
		conf.GetPort(),
		conf.GetUsername(),
		conf.GetName(),
		conf.GetPassword(),
	)

	if !conf.IsSSLMode() {
		connStr = fmt.Sprintf("%s sslmode=disable", connStr)
	}

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.Exec("CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;").Error
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&User{},
		&Check{},
		&Payload{},
		&Page{},
		&Incident{},
		&Metric{},
		&PageTeam{},
	)
	if err != nil {
		return nil, err
	}

	err = db.Exec("CREATE INDEX ON metrics (check_id, start_time DESC);").Error
	if err != nil {
		return nil, err
	}

	err = db.Exec(
		"SELECT create_hypertable('metrics', 'start_time', if_not_exists => TRUE, create_default_indexes => FALSE);",
	).Error
	if err != nil {
		return nil, err
	}

	return &Conn{db: db}, nil
}
