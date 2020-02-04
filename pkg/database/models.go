package database

import (
	"time"

	"github.com/jinzhu/gorm"
)

// User model.
type User struct {
	gorm.Model
	Email string `gorm:"UNIQUE;NOT NULL"`
	Name  string `gorm:"NOT NULL"`
}

// Check model.
type Check struct {
	gorm.Model
	OwnerID     uint
	Owner       User
	Interval    int       `gorm:"DEFAULT:30"`
	Timeout     int       `gorm:"DEFAULT:30"`
	InputType   string    `gorm:"NOT NULL"`
	InputValue  string    `gorm:"NOT NULL"`
	OutputType  string    `gorm:"NOT NULL"`
	OutputValue string    `gorm:"NOT NULL"`
	TargetType  string    `gorm:"NOT NULL"`
	TargetValue string    `gorm:"NOT NULL"`
	Title       string    `gorm:"NOT NULL"`
	Payloads    []Payload `gorm:"foreignkey:CheckID"`
}

// Payload model.
type Payload struct {
	gorm.Model
	CheckID uint
	Check   Check
	Type    string `gorm:"NOT NULL"`
	Value   string `gorm:"NOT NULL;TYPE:text"`
}

// Page model.
type Page struct {
	gorm.Model
	OwnerID     uint
	Owner       User
	Visibility  bool       `gorm:"DEFAULT:false"`
	Title       string     `gorm:"NOT NULL"`
	Description string     `gorm:"TYPE:text"`
	Checks      []Check    `gorm:"many2many:page_checks"`
	Team        []User     `gorm:"many2many:page_team"`
	Incidents   []Incident `gorm:"foreignkey:PageID"`
}

// Incident model.
type Incident struct {
	gorm.Model
	PageID      uint
	Page        Page
	TimeStamp   *time.Time `gorm:"NOT NULL"`
	Duration    int        `gorm:"NOT NULL"`
	Title       string     `gorm:"NOT NULL"`
	Description string     `gorm:"TYPE:text"`
	Resolved    bool       `gorm:"DEFAULT:false"`
}
