package database

import (
	"context"
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
func NewConn(ctx context.Context, conf Config) (*Conn, error) {
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

	err = db.WithContext(ctx).Exec("CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;").Error
	if err != nil {
		return nil, err
	}

	err = db.WithContext(ctx).AutoMigrate(
		&User{},
		&Check{},
		&Payload{},
		&Page{},
		&Incident{},
		&PageTeam{},
	)
	if err != nil {
		return nil, err
	}

	return &Conn{db: db}, nil
}
