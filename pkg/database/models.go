// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package database

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Bool for storing in models as a rune.
//
// 't', 'T', '1' and 1 represnt true, rest false.
type Bool rune

// True and False.
const (
	True  Bool = 't'
	False Bool = 'f'
)

// T tells if b is true or false.
func (b Bool) T() bool {
	return (b == 't' || b == 'T' || b == '1' || b == 1)
}

// User model.
type User struct {
	gorm.Model

	Email string `gorm:"UNIQUE;NOT NULL"`
	Name  string `gorm:"NOT NULL"`

	Checks    []Check    `gorm:"foreignkey:OwnerID"`
	Payloads  []Payload  `gorm:"foreignkey:OwnerID"`
	Pages     []Page     `gorm:"foreignkey:OwnerID"`
	TeamPages []PageTeam `gorm:"foreignkey:UserID"`
	Incidents []Incident `gorm:"foreignkey:OwnerID"`
}

// Check model.
type Check struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	OwnerID uint
	Owner   User

	Title string `gorm:"NOT NULL"`

	Interval time.Duration `gorm:"DEFAULT:30"`
	Timeout  time.Duration `gorm:"DEFAULT:30"`

	InputType  string `gorm:"NOT NULL"`
	InputValue string `gorm:"NOT NULL"`

	OutputType  string `gorm:"NOT NULL"`
	OutputValue string `gorm:"NOT NULL"`

	TargetType  string `gorm:"NOT NULL"`
	TargetValue string `gorm:"NOT NULL"`

	Payloads []Payload `gorm:"foreignkey:CheckID"`
	Metrics  []Metric  `gorm:"foreignkey:CheckID"`
}

// Payload model.
type Payload struct {
	gorm.Model

	Owner   User
	OwnerID uint

	CheckID string
	Check   Check

	Type  string `gorm:"NOT NULL"`
	Value string `gorm:"NOT NULL;TYPE:text"`
}

// Page model.
type Page struct {
	gorm.Model

	OwnerID uint
	Owner   User

	Title       string `gorm:"NOT NULL"`
	Description string `gorm:"TYPE:text"`
	Visibility  Bool   `gorm:"DEFAULT:102;size:256"`

	Checks    []Check    `gorm:"many2many:page_checks"`
	Incidents []Incident `gorm:"foreignkey:PageID"`
	Team      []PageTeam
}

// Incident model.
type Incident struct {
	gorm.Model

	Owner   User
	OwnerID uint

	PageID uint
	Page   Page

	Title       string `gorm:"NOT NULL"`
	Description string `gorm:"TYPE:text"`
	Resolved    Bool   `gorm:"DEFAULT:102;size:256"`

	TimeStamp time.Time     `gorm:"NOT NULL"`
	Duration  time.Duration `gorm:"NOT NULL"`
}

// Metric model.
type Metric struct {
	CheckID string
	Check   Check

	StartTime time.Time     `gorm:"NOT NULL"`
	Duration  time.Duration `gorm:"NOT NULL"`
	Timeout   bool          `gorm:"NOT NULL"`
	Success   bool          `gorm:"NOT NULL"`
}

// Various roles of a team member.
const (
	RoleDefault    = "DEFAULT"
	RoleMaintainer = "MAINTAINER"
	RoleAdmin      = "ADMIN"
)

// PageTeam model.
type PageTeam struct {
	Page   Page
	PageID uint `gorm:"primary_key;auto_increment:false"`

	User   User
	UserID uint `gorm:"primary_key;auto_increment:false"`

	Role string `gorm:"NOT NULL"`
}
