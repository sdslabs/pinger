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
	Interval    time.Duration `gorm:"DEFAULT:30"`
	Timeout     time.Duration `gorm:"DEFAULT:30"`
	InputType   string        `gorm:"NOT NULL"`
	InputValue  string        `gorm:"NOT NULL"`
	OutputType  string        `gorm:"NOT NULL"`
	OutputValue string        `gorm:"NOT NULL"`
	TargetType  string        `gorm:"NOT NULL"`
	TargetValue string        `gorm:"NOT NULL"`
	Title       string        `gorm:"NOT NULL"`
	Payloads    []Payload     `gorm:"foreignkey:CheckID"`
}

// PageTeam model.
type PageTeam struct {
	Page   *Page
	PageID int
	User   *User
	UserID int
	Role   string
}

// Payload model.
type Payload struct {
	gorm.Model
	Owner   User
	OwnerID uint
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
	Incidents   []Incident `gorm:"foreignkey:PageID"`
	Team        []PageTeam
}

// Incident model.
type Incident struct {
	gorm.Model
	Owner       User
	OwnerID     uint
	PageID      uint
	Page        Page
	TimeStamp   *time.Time    `gorm:"NOT NULL"`
	Duration    time.Duration `gorm:"NOT NULL"`
	Title       string        `gorm:"NOT NULL"`
	Description string        `gorm:"TYPE:text"`
	Resolved    bool          `gorm:"DEFAULT:false"`
}

// Metric model.
type Metric struct {
	CheckID   uint
	Check     Check
	StartTime *time.Time    `gorm:"NOT NULL"`
	Duration  time.Duration `gorm:"NOT NULL"`
	Timeout   bool          `gorm:"NOT NULL"`
	Success   bool          `gorm:"NOT NULL"`
}
