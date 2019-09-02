package database

import "github.com/jinzhu/gorm"

// SQLDB is an interface which implements api methods
type SQLDB interface {
	GetUserByID(id int) (User, error)
	GetUserByEmail(email string) (User, error)
	CreateUser(email, name string) error
	DeleteUserByID(id int) error
	DeleteUserByEmail(email string) error
	GetAllChecksByOwner(ownerID int) ([]Check, error)
	GetCheckByID(id int) (Check, error)
}

type sqldb struct {
	*gorm.DB
}

// GetUserByID gets user by id
func (db *sqldb) GetUserByID(id int) (User, error) {
	user := User{}
	tx := db.Where("id = ?", id).Find(&user)
	if tx.RecordNotFound() {
		return User{}, nil
	}
	return user, tx.Error
}

// GetUserByEmail gets user by email
func (db *sqldb) GetUserByEmail(email string) (User, error) {
	user := User{}
	tx := db.Where("email = ?", email).Find(&user)
	if tx.RecordNotFound() {
		return User{}, nil
	}
	return user, tx.Error
}

// CreateUser adds an entry for new user
func (db *sqldb) CreateUser(email, name string) error {
	user := User{
		Email: email,
		Name:  name,
	}
	tx := db.FirstOrCreate(&user, User{Email : email})
	return tx.Error
}

// DeleteUserByID deletes a user entry
func (db *sqldb) DeleteUserByID(id int) error {
	tx := db.Where("id = ?", id).Delete(&User{})
	return tx.Error
}

// DeleteUserByEmail deletes a user entry
func (db *sqldb) DeleteUserByEmail(email string) error {
	tx := db.Where("email = ?", email).Delete(&User{})
	return tx.Error
}

// GetAllChecksByOwner gets all the checks in owned by the user
func (db *sqldb) GetAllChecksByOwner(ownerID int) ([]Check, error) {
	checks := []Check{}
	tx := db.Where("owner_id = ?", ownerID).Preload("Payloads").Preload("Owner").Find(&checks)
	if tx.RecordNotFound() {
		return nil, nil
	}
	return checks, tx.Error
}

// GetCheckByID gets a check by its ID
func (db *sqldb) GetCheckByID(id int) (Check, error) {
	check := Check{}
	tx := db.Where("id = ?", id).Preload("Payloads").Preload("Owner").Find(&check)
	if tx.RecordNotFound() {
		return Check{}, nil
	}
	return check, tx.Error
}
