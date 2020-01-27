package database

import (
	"time"

	"github.com/jinzhu/gorm"
)

// SQLDB is an interface which implements api methods
type SQLDB interface {
	GetUserByID(id int) (User, error)
	GetUserByEmail(email string) (User, error)
	CreateUser(email, name string) (User, error)
	UpdateUserNameByID(id uint, name string) (User, error)
	UpdateUserNameByEmail(email, name string) (User, error)
	DeleteUserByID(id int) error
	DeleteUserByEmail(email string) error
	GetAllChecksByOwner(ownerID int) ([]Check, error)
	GetCheckByID(id int) (Check, error)
	CreateCheck(
		ownerID, interval, timeout int,
		inputType, inputValue, outputType, outputValue, targetType, targetValue, title string,
		payloads []Payload) (Check, error)
	UpdateCheckByID(id uint, check *Check) (Check, error)
	DeleteCheckByID(id int) error
	GetAllPayloadsByCheck(checkID int) ([]Payload, error)
	GetPayloadByID(id int) (Payload, error)
	CreatePayload(checkID int, payloadType, value string) (Payload, error)
	UpdatePayloadByID(id uint, payload *Payload) (Payload, error)
	DeletePayloadByID(id int) error
	AddPayloadsToCheck(checkID uint, payloads []*Payload) error
	// RemovePayloadsFromCheck(checkID uint, payloads []*Payload) error
	GetAllPagesByOwner(ownerID int) ([]Page, error)
	GetPageByID(id int) (Page, error)
	CreatePage(ownerID int, visibility bool, title, description string, incidents []Incident) (Page, error)
	UpdatePageByID(id uint, page *Page) (Page, error)
	DeletePageByID(id int) error
	GetAllIncidentsByPage(pageID int) ([]Incident, error)
	GetIncidentByID(id int) (Incident, error)
	CreateIncident(pageID int, timestamp *time.Time, duration int, title, description string, resolved bool) (Incident, error)
	UpdateIncidentByID(id uint, incident *Incident) (Incident, error)
	DeleteIncidentByID(id int) error
	AddIncidentsToPage(pageID uint, incidents []*Incident) error
	// RemoveIncidentsFromPage(pageID uint, incidents []*Incident) error
	AddChecksToPage(pageID uint, checks []*Check) error
	RemoveChecksFromPage(pageID uint, checks []*Check) error
	AddMembersToPageTeam(pageID uint, users []*User) error
	RemoveMembersFromPageTeam(pageID uint, users []*User) error
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
// If the user already exists, does nothing
func (db *sqldb) CreateUser(email, name string) (User, error) {
	u, err := db.GetUserByEmail(email)
	if err != nil {
		return User{}, err
	}
	if u.Email == email {
		return u, nil
	}
	user := User{
		Email: email,
		Name:  name,
	}
	tx := db.Create(&user)
	return user, tx.Error
}

// UpdateUserNameByID updates the Name of the user
func (db *sqldb) UpdateUserNameByID(id uint, name string) (User, error) {
	user := User{}
	user.ID = id
	tx := db.Model(&user).Update("name", name)
	return user, tx.Error
}

// UpdateUserNameByEmail updates the Name of the user
func (db *sqldb) UpdateUserNameByEmail(email, name string) (User, error) {
	user := User{Email: email}
	tx := db.Model(&user).Update("name", name)
	return user, tx.Error
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

// CreateCheck creates a new check
func (db *sqldb) CreateCheck(
	ownerID, interval, timeout int,
	inputType, inputValue, outputType, outputValue, targetType, targetValue, title string,
	payloads []Payload) (Check, error) {
	check := Check{
		OwnerID:     ownerID,
		Interval:    interval,
		Timeout:     timeout,
		InputType:   inputType,
		InputValue:  inputValue,
		OutputType:  outputType,
		OutputValue: outputValue,
		TargetType:  targetType,
		TargetValue: targetValue,
		Title:       title,
		Payloads:    payloads,
	}
	tx := db.Create(&check)
	return check, tx.Error
}

// UpdateCheckByID updates the check for given ID
func (db *sqldb) UpdateCheckByID(id uint, check *Check) (Check, error) {
	c := Check{}
	c.ID = id
	tx := db.Model(&c).Updates(*check)
	return c, tx.Error
}

// DeleteCheckByID deletes check corresponding to given ID
func (db *sqldb) DeleteCheckByID(id int) error {
	tx := db.Where("id = ?", id).Unscoped().Delete(&Check{})
	return tx.Error
}

// GetAllPayloadsByCheck gets all the payloads belonging to a check
func (db *sqldb) GetAllPayloadsByCheck(checkID int) ([]Payload, error) {
	payloads := []Payload{}
	tx := db.Where("check_id = ?", checkID).Preload("Check").Find(&payloads)
	if tx.RecordNotFound() {
		return nil, nil
	}
	return payloads, tx.Error
}

// GetPayloadByID gets a payload corresponding to the ID
func (db *sqldb) GetPayloadByID(id int) (Payload, error) {
	payload := Payload{}
	tx := db.Where("id = ?", id).Preload("Check").Find(&payload)
	if tx.RecordNotFound() {
		return Payload{}, nil
	}
	return payload, tx.Error
}

// CreatePayload creates a payload with given type and value
func (db *sqldb) CreatePayload(checkID int, payloadType, value string) (Payload, error) {
	payload := Payload{
		CheckID: checkID,
		Type:    payloadType,
		Value:   value,
	}
	tx := db.Create(&payload)
	return payload, tx.Error
}

// UpdatePayloadByID updates the payload for given ID
func (db *sqldb) UpdatePayloadByID(id uint, payload *Payload) (Payload, error) {
	p := Payload{}
	p.ID = id
	tx := db.Model(&p).Updates(*payload)
	return p, tx.Error
}

// DeletePayloadByID deletes a payload corresponding to given ID
func (db *sqldb) DeletePayloadByID(id int) error {
	tx := db.Where("id = ?", id).Unscoped().Delete(&Payload{})
	return tx.Error
}

// AddPayloadsToPage adds multiple payloads to page
func (db *sqldb) AddPayloadsToCheck(checkID uint, payloads []*Payload) error {
	check := Check{}
	check.ID = checkID
	tx := db.Model(&check).Association("Payloads").Append(payloads)
	return tx.Error
}

// ***
// [TODO]
// *gorm.Association.Delete() only deletes the relationship but not the elements
// currently not using this for bulk delete
// ***
// // RemovePayloadsFromCheck adds multiple checks to page
// func (db *sqldb) RemovePayloadsFromCheck(checkID uint, payloads []*Payload) error {
// 	check := Check{}
// 	check.ID = checkID
// 	tx := db.Model(&check).Association("Payloads").Delete(payloads)
// 	return tx.Error
// }

// GetAllPagesByOwner gets all the pages in owned by the user
func (db *sqldb) GetAllPagesByOwner(ownerID int) ([]Page, error) {
	pages := []Page{}
	tx := db.Where("owner_id = ?", ownerID).Preload("Checks").Preload("Team").Preload("Incidents").Preload("Owner").Find(&pages)
	if tx.RecordNotFound() {
		return nil, nil
	}
	return pages, tx.Error
}

// GetPageByID gets a check by its ID
func (db *sqldb) GetPageByID(id int) (Page, error) {
	page := Page{}
	tx := db.Where("id = ?", id).Preload("Checks").Preload("Team").Preload("Incidents").Preload("Owner").Find(&page)
	if tx.RecordNotFound() {
		return Page{}, nil
	}
	return page, tx.Error
}

// CreatePage creates a new page
func (db *sqldb) CreatePage(ownerID int, visibility bool, title, description string, incidents []Incident) (Page, error) {
	page := Page{
		OwnerID:     ownerID,
		Visibility:  visibility,
		Title:       title,
		Description: description,
		Incidents:   incidents,
	}
	tx := db.Create(&page)
	return page, tx.Error
}

// UpdatePageByID updates the page for given ID
func (db *sqldb) UpdatePageByID(id uint, page *Page) (Page, error) {
	p := Page{}
	p.ID = id
	tx := db.Model(&p).Updates(*page)
	return p, tx.Error
}

// DeletePageByID deletes a page corresponding to the given ID
func (db *sqldb) DeletePageByID(id int) error {
	tx := db.Where("id = ?", id).Delete(&Page{})
	return tx.Error
}

// GetAllIncidentsByPage gets all the incidents for the given page ID
func (db *sqldb) GetAllIncidentsByPage(pageID int) ([]Incident, error) {
	incidents := []Incident{}
	tx := db.Where("page_id = ?", pageID).Preload("Page").Find(&incidents)
	if tx.RecordNotFound() {
		return nil, nil
	}
	return incidents, tx.Error
}

// GetIncidentByID gets incident corresponding to given ID
func (db *sqldb) GetIncidentByID(id int) (Incident, error) {
	incident := Incident{}
	tx := db.Where("id = ?", id).Preload("Page").Find(&incident)
	if tx.RecordNotFound() {
		return Incident{}, nil
	}
	return incident, tx.Error
}

// CreateIncident creates an incident with given type and value
func (db *sqldb) CreateIncident(
	pageID int,
	timestamp *time.Time,
	duration int,
	title, description string,
	resolved bool) (Incident, error) {
	incident := Incident{
		PageID:      pageID,
		TimeStamp:   timestamp,
		Duration:    duration,
		Title:       title,
		Description: description,
		Resolved:    resolved,
	}
	tx := db.Create(&incident)
	return incident, tx.Error
}

// UpdateIncidentByID updates the incident for given ID
func (db *sqldb) UpdateIncidentByID(id uint, incident *Incident) (Incident, error) {
	i := Incident{}
	i.ID = id
	tx := db.Model(&i).Updates(*incident)
	return i, tx.Error
}

// DeleteIncidentByID deletes a Incident corresponding to given ID
func (db *sqldb) DeleteIncidentByID(id int) error {
	tx := db.Where("id = ?", id).Unscoped().Delete(&Incident{})
	return tx.Error
}

// AddIncidentsToPage adds multiple incidents to page
func (db *sqldb) AddIncidentsToPage(pageID uint, incidents []*Incident) error {
	page := Page{}
	page.ID = pageID
	tx := db.Model(&page).Association("Incidents").Append(incidents)
	return tx.Error
}

// ***
// [TODO]
// *gorm.Association.Delete() only deletes the relationship but not the elements
// currently not using this for bulk delete
// ***
// // RemoveIncidentsFromPage adds multiple checks to page
// func (db *sqldb) RemoveIncidentsFromPage(pageID uint, incidents []*Incident) error {
// 	page := Page{}
// 	page.ID = pageID
// 	tx := db.Model(&page).Association("Incidents").Delete(incidents)
// 	return tx.Error
// }

// AddChecksToPage adds multiple checks to page
func (db *sqldb) AddChecksToPage(pageID uint, checks []*Check) error {
	page := Page{}
	page.ID = pageID
	tx := db.Model(&page).Association("Checks").Append(checks)
	return tx.Error
}

// RemoveChecksFromPage adds multiple checks to page
func (db *sqldb) RemoveChecksFromPage(pageID uint, checks []*Check) error {
	page := Page{}
	page.ID = pageID
	tx := db.Model(&page).Association("Checks").Delete(checks)
	return tx.Error
}

// AddMembersToPageTeam adds multiple checks to page
func (db *sqldb) AddMembersToPageTeam(pageID uint, users []*User) error {
	page := Page{}
	page.ID = pageID
	tx := db.Model(&page).Association("Team").Append(users)
	return tx.Error
}

// RemoveMembersFromPageTeam adds multiple checks to page
func (db *sqldb) RemoveMembersFromPageTeam(pageID uint, users []*User) error {
	page := Page{}
	page.ID = pageID
	tx := db.Model(&page).Association("Team").Delete(users)
	return tx.Error
}
