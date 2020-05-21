package database

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrRecordNotFound is the error returned when the database has no matching record.
var ErrRecordNotFound = errors.New("record not found")

// GetUserByID gets user by ID.
func GetUserByID(id uint) (*User, error) {
	user := User{}
	tx := db.Where("id = ?", id).Find(&user)
	if tx.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	return &user, tx.Error
}

// GetUserByEmail gets user by Email.
func GetUserByEmail(email string) (*User, error) {
	user := User{}
	tx := db.Where("email = ?", email).Find(&user)
	if tx.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	return &user, tx.Error
}

// CreateUser adds an entry for new user. If the user already exists, does nothing.
//
// BUG(murex971): once deleted, a user could not be created with same email.
func CreateUser(user *User) (*User, error) {
	u, err := GetUserByEmail(user.Email)
	if err != nil && err != ErrRecordNotFound {
		return nil, err
	}
	if err == nil && u.Email == user.Email {
		return u, nil
	}
	tx := db.Create(user)
	return user, tx.Error
}

// UpdateUserByID updates the user for given ID.
func UpdateUserByID(id uint, user *User) (*User, error) {
	u := User{}
	u.ID = id
	tx := db.Model(&u).Updates(*user)
	if tx.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	return &u, tx.Error
}

// UpdateUserByEmail updates the user for given email.
func UpdateUserByEmail(email string, user *User) (*User, error) {
	u := User{}
	u.Email = email
	tx := db.Model(&u).Updates(*user)
	if tx.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	return &u, tx.Error
}

// DeleteUserByID deletes a user entry.
func DeleteUserByID(id uint) error {
	tx := db.Where("id = ?", id).Delete(&User{})
	if tx.RecordNotFound() {
		return ErrRecordNotFound
	}
	return tx.Error
}

// DeleteUserByEmail deletes a user entry.
func DeleteUserByEmail(email string) error {
	tx := db.Where("email = ?", email).Delete(&User{})
	if tx.RecordNotFound() {
		return ErrRecordNotFound
	}
	return tx.Error
}

// GetAllChecksByOwner gets all the checks in owned by the user.
func GetAllChecksByOwner(ownerID uint) ([]Check, error) {
	checks := []Check{}
	tx := db.Where("owner_id = ?", ownerID).Preload("Payloads").Preload("Owner").Find(&checks)
	if tx.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	return checks, tx.Error
}

// GetCheckByID gets a check by its ID.
func GetCheckByID(id uint) (*Check, error) {
	check := Check{}
	tx := db.Where("id = ?", id).Preload("Payloads").Preload("Owner").Find(&check)
	if tx.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	return &check, tx.Error
}

// CreateCheck creates a new check.
func CreateCheck(check *Check) (*Check, error) {
	tx := db.Create(check)
	return check, tx.Error
}

// UpdateCheckByID updates the check for given ID.
func UpdateCheckByID(id, ownerID uint, check *Check) (*Check, error) {
	c := Check{}
	c.ID = id
	tx := db.Model(&c).Where("owner_id = ?", ownerID).Updates(*check)
	if tx.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	return &c, tx.Error
}

// DeleteCheckByID deletes check corresponding to given ID.
func DeleteCheckByID(id, ownerID uint) error {
	tx := db.Where("id = ? AND owner_id = ?", id, ownerID).Unscoped().Delete(&Check{})
	if tx.RecordNotFound() {
		return ErrRecordNotFound
	}
	return tx.Error
}

// GetAllPayloadsByCheck gets all the payloads belonging to a check.
func GetAllPayloadsByCheck(checkID uint) ([]Payload, error) {
	payloads := []Payload{}
	tx := db.Where("check_id = ?", checkID).Preload("Check").Preload("Owner").Find(&payloads)
	if tx.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	return payloads, tx.Error
}

// GetPayloadByID gets a payload corresponding to the ID.
func GetPayloadByID(id uint) (*Payload, error) {
	payload := Payload{}
	tx := db.Where("id = ?", id).Preload("Check").Preload("Owner").Find(&payload)
	if tx.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	return &payload, tx.Error
}

// CreatePayload creates a payload with given type and value.
func CreatePayload(payload *Payload) (*Payload, error) {
	tx := db.Create(payload)
	return payload, tx.Error
}

// UpdatePayloadByID updates the payload for given ID.
func UpdatePayloadByID(id, ownerID, checkID uint, payload *Payload) (*Payload, error) {
	p := Payload{}
	p.ID = id
	tx := db.Model(&p).Where("check_id = ? AND owner_id = ?", checkID, ownerID).Updates(*payload)
	if tx.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	return &p, tx.Error
}

// DeletePayloadByID deletes a payload corresponding to given ID.
func DeletePayloadByID(id, ownerID, checkID uint) error {
	tx := db.Where("id = ? AND check_id = ? AND owner_id = ?", id, checkID, ownerID).Unscoped().Delete(&Payload{})
	if tx.RecordNotFound() {
		return ErrRecordNotFound
	}
	return tx.Error
}

// AddPayloadsToCheck adds multiple payloads to page.
func AddPayloadsToCheck(ownerID, checkID uint, payloads []*Payload) error {
	check := Check{}
	check.ID = checkID
	return db.Model(&check).Where("owner_id = ?", ownerID).Association("Payloads").Append(payloads).Error
}

// RemovePayloadsFromCheck removes multiple checks from a page.
//
// BUG(vrongmeal): This only removes the relationship and not the payloads.
func RemovePayloadsFromCheck(ownerID, checkID uint, payloads []*Payload) error {
	check := Check{}
	check.ID = checkID
	return db.Model(&check).Where("owner_id = ?", ownerID).Association("Payloads").Delete(payloads).Error
}

// GetAllPagesByOwner gets all the pages in owned by the user.
func GetAllPagesByOwner(ownerID uint) ([]Page, error) {
	pages := []Page{}
	tx := db.Where("owner_id = ?", ownerID).
		Preload("Checks").
		Preload("Team.User").
		Preload("Incidents").
		Preload("Owner").
		Find(&pages)
	if tx.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	return pages, tx.Error
}

// GetPageByID gets a check by its ID.
func GetPageByID(id uint) (*Page, error) {
	page := Page{}
	tx := db.Where("id = ?", id).
		Preload("Checks").
		Preload("Team.User").
		Preload("Incidents").
		Preload("Owner").
		Find(&page)
	if tx.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	return &page, tx.Error
}

// CreatePage creates a new page.
func CreatePage(page *Page) (*Page, error) {
	tx := db.Create(page)
	return page, tx.Error
}

// UpdatePageByID updates the page for given ID.
func UpdatePageByID(id, ownerID uint, page *Page) (*Page, error) {
	p := Page{}
	p.ID = id
	tx := db.Model(&p).Where("owner_id = ?", ownerID).Updates(*page)
	if tx.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	return &p, tx.Error
}

// DeletePageByID deletes a page corresponding to the given ID.
func DeletePageByID(id, ownerID uint) error {
	tx := db.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&Page{})
	if tx.RecordNotFound() {
		return ErrRecordNotFound
	}
	return tx.Error
}

// GetAllIncidentsByPage gets all the incidents for the given page ID.
func GetAllIncidentsByPage(pageID uint) ([]Incident, error) {
	incidents := []Incident{}
	tx := db.Where("page_id = ?", pageID).Preload("Page").Preload("Owner").Find(&incidents)
	if tx.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	return incidents, tx.Error
}

// GetIncidentByID gets incident corresponding to given ID.
func GetIncidentByID(id uint) (*Incident, error) {
	incident := Incident{}
	tx := db.Where("id = ?", id).Preload("Page").Preload("Owner").Find(&incident)
	if tx.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	return &incident, tx.Error
}

// CreateIncident creates an incident with given type and value.
func CreateIncident(incident *Incident) (*Incident, error) {
	tx := db.Create(incident)
	return incident, tx.Error
}

// UpdateIncidentByID updates the incident for given ID.
func UpdateIncidentByID(id, ownerID, pageID uint, incident *Incident) (*Incident, error) {
	i := Incident{}
	i.ID = id
	tx := db.Model(&i).Where("page_id = ? AND owner_id = ?", pageID, ownerID).Updates(*incident)
	if tx.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	return &i, tx.Error
}

// DeleteIncidentByID deletes a Incident corresponding to given ID.
func DeleteIncidentByID(id, ownerID, pageID uint) error {
	tx := db.Where("id = ? AND page_id = ? AND owner_id = ?", id, pageID, ownerID).Unscoped().Delete(&Incident{})
	if tx.RecordNotFound() {
		return ErrRecordNotFound
	}
	return tx.Error
}

// AddIncidentsToPage adds multiple incidents to page.
func AddIncidentsToPage(ownerID, pageID uint, incidents []*Incident) error {
	page := Page{}
	page.ID = pageID
	return db.Model(&page).Where("owner_id = ?", ownerID).Association("Incidents").Append(incidents).Error
}

// RemoveIncidentsFromPage adds multiple checks to page.
//
// BUG(vrongmeal): This only removes the relationship and not the incidents.
func RemoveIncidentsFromPage(ownerID, pageID uint, incidents []*Incident) error {
	page := Page{}
	page.ID = pageID
	return db.Model(&page).Where("owner_id = ?", ownerID).Association("Incidents").Delete(incidents).Error
}

// AddChecksToPage adds multiple checks to page.
func AddChecksToPage(ownerID, pageID uint, checks []*Check) error {
	page := Page{}
	page.ID = pageID
	return db.Model(&page).Where("owner_id = ?", ownerID).Association("Checks").Append(checks).Error
}

// RemoveChecksFromPage removes multiple checks from page.
func RemoveChecksFromPage(ownerID, pageID uint, checks []*Check) error {
	page := Page{}
	page.ID = pageID
	return db.Model(&page).Where("owner_id = ?", ownerID).Association("Checks").Delete(checks).Error
}

// AddMembersToPageTeam adds users as new members to a team.
func AddMembersToPageTeam(ownerID, pageID uint, users []*User) error {
	page := Page{}
	page.ID = pageID
	return db.Model(&page).Where("owner_id = ?", ownerID).Association("Team.User").Append(users).Error
}

// RemoveMembersFromPageTeam removes members from a team.
func RemoveMembersFromPageTeam(ownerID, pageID uint, users []*User) error {
	page := Page{}
	page.ID = pageID
	return db.Model(&page).Where("owner_id = ?", ownerID).Association("Team.User").Delete(users).Error
}

// GetMetricsByCheckAndStartTime fetches metrics from the metrics hypertable for the given check ID.
// It accepts a `startTime` parameter that fetches metrics for the check from given time.
func GetMetricsByCheckAndStartTime(checkID uint, startTime time.Time) ([]Metric, error) {
	metrics := []Metric{}
	tx := db.Where("check_id = ? AND start_time > ?", checkID, startTime).Order("start_time DESC").Find(&metrics)
	if tx.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	return metrics, tx.Error
}

// GetMetricsByCheckAndDuration fetches metrics from the metrics hypertable for the given check ID.
// It accepts a `duration` parameter that fetches metrics for the check in the past `duration time.Duration`.
func GetMetricsByCheckAndDuration(checkID uint, duration time.Duration) ([]Metric, error) {
	startTime := time.Now().Add(-1 * duration)
	return GetMetricsByCheckAndStartTime(checkID, startTime)
}

// GetMetricsByPageAndStartTime fetches metrics for all the checks in a page for the given start time.
func GetMetricsByPageAndStartTime(pageID uint, startTime time.Time) ([]Metric, error) {
	checkIDs := []uint{}
	tx1 := db.Table("page_checks").Where("page_id = ?", pageID).Pluck("check_id", &checkIDs)
	if tx1.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	if err := tx1.Error; err != nil {
		return nil, err
	}
	metrics := []Metric{}
	tx2 := db.Where("check_id IN (?) AND start_time > ?", checkIDs, startTime).
		Order("start_time DESC").
		Find(&metrics)
	if tx2.RecordNotFound() {
		return nil, ErrRecordNotFound
	}
	return metrics, tx2.Error
}

// GetMetricsByPageAndDuration fetches metrics for all the checks in a page for the given duration.
func GetMetricsByPageAndDuration(pageID uint, duration time.Duration) ([]Metric, error) {
	startTime := time.Now().Add(-1 * duration)
	return GetMetricsByPageAndStartTime(pageID, startTime)
}

// CreateMetrics inserts multiple metrics into TimescaleDB Hypertable.
func CreateMetrics(metrics []Metric) error {
	// We build a raw query since Gorm doesn't support bulk insert.
	// Since there is no `string` and there is no user input we can
	// safely build the raw query without worrying about injection.
	if len(metrics) == 0 {
		return nil
	}
	q := "INSERT INTO metrics (check_id, start_time, duration, timeout, success) VALUES %s;"
	timeFormat := "2006-01-02 15:04:05.000000-07:00" // Supported by PostgreSQL
	vals := []string{}
	for i := range metrics {
		val := fmt.Sprintf("(%d, '%s', %d, %t, %t)",
			metrics[i].CheckID,
			metrics[i].StartTime.Format(timeFormat),
			metrics[i].Duration,
			metrics[i].Timeout,
			metrics[i].Success)
		vals = append(vals, val)
	}
	args := strings.Join(vals, ", ")
	query := fmt.Sprintf(q, args)
	return db.Exec(query).Error
}
