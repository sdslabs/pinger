// Copyright (c) 2020 SDSLabs
// Use of this source code is governed by an MIT license
// details of which can be found in the LICENSE file.

package database

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	// ErrRecordNotFound is the error returned when the database has no matching
	// record.
	//
	// Since gorm v2 (1.20), it returns an error rather than a check using method,
	// so this error is just an alias of same error which can be used to check if
	// the record did not exist.
	ErrRecordNotFound = gorm.ErrRecordNotFound

	// ErrNilPointer is the error returned in case when the pointer is nil. This
	// can be checked using `errors.Unwrap`.
	ErrNilPointer = errors.New("pointer cannot be nil")
)

// rawUserWithID returns an empty user with the given ID.
func rawUserWithID(id uint) User {
	user := User{}
	user.ID = id
	return user
}

// rawUserWithID returns an empty user with the given email.
func rawUserWithEmail(email string) User {
	user := User{}
	user.Email = email
	return user
}

// CreateUser creates a new user in the database. It simply returns the user
// with the same email if the user exists.
func (c *Conn) CreateUser(ctx context.Context, user *User) (*User, error) {
	if user == nil {
		return nil, fmt.Errorf("*User: %w", ErrNilPointer)
	}

	u, err := c.GetUserByEmail(ctx, user.Email, GetUserOpts{})
	if err != nil && err != ErrRecordNotFound {
		return nil, err
	}

	if err == nil && u.Email == user.Email {
		return u, nil
	}

	err = c.db.WithContext(ctx).Create(user).Error
	return user, err
}

// GetUserOpts are the get options for user relations. Objects are preloaded
// using the options set here.
type GetUserOpts struct {
	Checks    bool
	Payloads  bool
	Pages     bool
	TeamPages bool
	Incidents bool
}

// getUser gets a user with the specified "where" condition.
func (c *Conn) getUser(ctx context.Context, where *User, opts GetUserOpts) (*User, error) {
	tx := c.db.WithContext(ctx).Where(where)

	if opts.Checks {
		tx = tx.Preload("Checks")
	}

	if opts.Payloads {
		tx = tx.Preload("Payloads")
	}

	if opts.Pages {
		tx = tx.Preload("Pages")
	}

	if opts.TeamPages {
		tx = tx.Preload("TeamPages.Page")
	}

	if opts.Incidents {
		tx = tx.Preload("Incidents")
	}

	user := User{}
	tx = tx.Find(&user)
	return &user, tx.Error
}

// GetUserByID gets user by ID.
func (c *Conn) GetUserByID(ctx context.Context, id uint, opts GetUserOpts) (*User, error) {
	u := rawUserWithID(id)
	return c.getUser(ctx, &u, opts)
}

// GetUserByEmail gets user by Email.
func (c *Conn) GetUserByEmail(ctx context.Context, email string, opts GetUserOpts) (*User, error) {
	u := rawUserWithEmail(email)
	return c.getUser(ctx, &u, opts)
}

// UpdateUserByID updates the user for given ID.
func (c *Conn) UpdateUserByID(ctx context.Context, id uint, user *User) (*User, error) {
	if user == nil {
		return nil, fmt.Errorf("*User: %w", ErrNilPointer)
	}

	u := rawUserWithID(id)

	tx := c.db.WithContext(ctx).Model(User{}).Where(&u).Updates(*user)
	return &u, tx.Error
}

// UpdateUserByEmail updates the user for given email.
func (c *Conn) UpdateUserByEmail(ctx context.Context, email string, user *User) (*User, error) {
	if user == nil {
		return nil, fmt.Errorf("*User: %w", ErrNilPointer)
	}

	u := rawUserWithEmail(email)

	tx := c.db.WithContext(ctx).Model(User{}).Where(&u).Updates(*user)
	return &u, tx.Error
}

// DeleteUserByID deletes a user entry.
func (c *Conn) DeleteUserByID(ctx context.Context, id uint) error {
	u := rawUserWithID(id)

	tx := c.db.WithContext(ctx).Where(&u).Unscoped().Delete(&User{})
	return tx.Error
}

// DeleteUserByEmail deletes a user entry.
func (c *Conn) DeleteUserByEmail(ctx context.Context, email string) error {
	u := rawUserWithEmail(email)

	tx := c.db.WithContext(ctx).Where(&u).Unscoped().Delete(&User{})
	return tx.Error
}

// rawCheckWithID returns an empty check with the given ID.
func rawCheckWithID(ownerID uint, checkID string) Check {
	check := Check{}
	check.OwnerID = ownerID
	check.ID = checkID
	return check
}

// CreateCheck creates a new check.
func (c *Conn) CreateCheck(ctx context.Context, ownerID uint, check *Check) (*Check, error) {
	if check == nil {
		return nil, fmt.Errorf("*Check: %w", ErrNilPointer)
	}

	check.OwnerID = ownerID
	err := c.db.WithContext(ctx).Create(check).Error
	return check, err
}

// GetCheckOpts are the options to preload check associations.
type GetCheckOpts struct {
	Owner bool

	Payloads bool
}

// GetCheck gets a check from given checkID.
func (c *Conn) GetCheck(ctx context.Context, ownerID uint, checkID string, opts GetCheckOpts) (*Check, error) {
	ch := rawCheckWithID(ownerID, checkID)
	tx := c.db.WithContext(ctx).Where(ch)

	if opts.Owner {
		tx = tx.Preload("Owner")
	}

	if opts.Payloads {
		tx = tx.Preload("Payloads")
	}

	check := Check{}
	tx = tx.Find(&check)
	return &check, tx.Error
}

// UpdateCheck updates a check with the given ID.
func (c *Conn) UpdateCheck(ctx context.Context, ownerID uint, checkID string, check *Check) (*Check, error) {
	if check == nil {
		return nil, fmt.Errorf("*Check: %w", ErrNilPointer)
	}

	ch := rawCheckWithID(ownerID, checkID)

	tx := c.db.WithContext(ctx).Model(Check{}).Where(&ch).Updates(*check)
	return &ch, tx.Error
}

// DeleteCheck deletes the check with the given ID.
func (c *Conn) DeleteCheck(ctx context.Context, ownerID uint, checkID string) error {
	ch := rawCheckWithID(ownerID, checkID)

	tx := c.db.WithContext(ctx).Where(&ch).Unscoped().Delete(&Check{})
	return tx.Error
}

// rawPayloadWithID returns an empty payload with the given ID.
func rawPayloadWithID(ownerID, payloadID uint, checkID string) Payload {
	payload := Payload{}
	payload.OwnerID = ownerID
	payload.CheckID = checkID
	payload.ID = payloadID
	return payload
}

// CreatePayload creates a new payload.
func (c *Conn) CreatePayload(ctx context.Context, ownerID uint, checkID string, payload *Payload) (*Payload, error) {
	if payload == nil {
		return nil, fmt.Errorf("*Payload: %w", ErrNilPointer)
	}

	payload.OwnerID = ownerID
	payload.CheckID = checkID
	err := c.db.WithContext(ctx).Create(payload).Error
	return payload, err
}

// GetPayloadOpts are the options to preload payload associations.
type GetPayloadOpts struct {
	Owner bool
	Check bool
}

// GetPayload gets a payload from given payloadID.
func (c *Conn) GetPayload(ctx context.Context, ownerID, payloadID uint, checkID string, opts GetPayloadOpts) (*Payload, error) {
	p := rawPayloadWithID(ownerID, payloadID, checkID)
	tx := c.db.WithContext(ctx).Where(p)

	if opts.Owner {
		tx = tx.Preload("Owner")
	}

	if opts.Check {
		tx = tx.Preload("Check")
	}

	payload := Payload{}
	tx = tx.Find(&payload)
	return &payload, tx.Error
}

// UpdatePayload updates a payload with the given ID.
func (c *Conn) UpdatePayload(ctx context.Context, ownerID, payloadID uint, checkID string, payload *Payload) (*Payload, error) {
	if payload == nil {
		return nil, fmt.Errorf("*Payload: %w", ErrNilPointer)
	}

	p := rawPayloadWithID(ownerID, payloadID, checkID)

	tx := c.db.WithContext(ctx).Model(Payload{}).Where(&p).Updates(*payload)
	return &p, tx.Error
}

// DeletePayload deletes the payload with the given ID.
func (c *Conn) DeletePayload(ctx context.Context, ownerID, payloadID uint, checkID string) error {
	p := rawPayloadWithID(ownerID, payloadID, checkID)

	tx := c.db.WithContext(ctx).Where(&p).Unscoped().Delete(&Payload{})
	return tx.Error
}

// rawPageWithID returns an empty page with the given ID.
func rawPageWithID(ownerID, pageID uint) Page {
	page := Page{}
	page.OwnerID = ownerID
	page.ID = pageID
	return page
}

// CreatePage creates a new page.
func (c *Conn) CreatePage(ctx context.Context, ownerID uint, page *Page) (*Page, error) {
	if page == nil {
		return nil, fmt.Errorf("*Page: %w", ErrNilPointer)
	}

	page.OwnerID = ownerID
	err := c.db.WithContext(ctx).Create(page).Error
	return page, err
}

// GetPageOpts are the options to preload page associations.
type GetPageOpts struct {
	Owner bool

	Checks    bool
	Incidents bool
	Team      bool
}

// GetPage gets a page from given pageID.
func (c *Conn) GetPage(ctx context.Context, ownerID, pageID uint, opts GetPageOpts) (*Page, error) {
	p := rawPageWithID(ownerID, pageID)
	tx := c.db.WithContext(ctx).Where(p)

	if opts.Owner {
		tx = tx.Preload("Owner")
	}

	if opts.Checks {
		tx = tx.Preload("Checks")
	}

	if opts.Incidents {
		tx = tx.Preload("Incidents")
	}

	if opts.Team {
		tx = tx.Preload("Team.User")
	}

	page := Page{}
	tx = tx.Find(&page)
	return &page, tx.Error
}

// UpdatePage updates a page with the given ID.
func (c *Conn) UpdatePage(ctx context.Context, ownerID, pageID uint, page *Page) (*Page, error) {
	if page == nil {
		return nil, fmt.Errorf("*Page: %w", ErrNilPointer)
	}

	p := rawPageWithID(ownerID, pageID)

	tx := c.db.WithContext(ctx).Model(Page{}).Where(&p).Updates(*page)
	return &p, tx.Error
}

// DeletePage deletes the page with the given ID.
func (c *Conn) DeletePage(ctx context.Context, ownerID, pageID uint) error {
	p := rawPageWithID(ownerID, pageID)

	tx := c.db.WithContext(ctx).Where(&p).Unscoped().Delete(&Page{})
	return tx.Error
}

// rawIncidentWithID returns an empty incident with the given ID.
func rawIncidentWithID(ownerID, pageID, incidentID uint) Incident {
	incident := Incident{}
	incident.OwnerID = ownerID
	incident.PageID = pageID
	incident.ID = incidentID
	return incident
}

// CreateIncident creates a new incident.
func (c *Conn) CreateIncident(ctx context.Context, ownerID, pageID uint, incident *Incident) (*Incident, error) {
	if incident == nil {
		return nil, fmt.Errorf("*Incident: %w", ErrNilPointer)
	}

	incident.OwnerID = ownerID
	incident.PageID = pageID
	err := c.db.WithContext(ctx).Create(incident).Error
	return incident, err
}

// GetIncidentOpts are the options to preload payload associations.
type GetIncidentOpts struct {
	Owner bool
	Page  bool
}

// GetIncident gets an incident from given incidentID.
func (c *Conn) GetIncident(ctx context.Context, ownerID, pageID, incidentID uint, opts GetIncidentOpts) (*Incident, error) {
	i := rawIncidentWithID(ownerID, pageID, incidentID)
	tx := c.db.WithContext(ctx).Where(i)

	if opts.Owner {
		tx = tx.Preload("Owner")
	}

	if opts.Page {
		tx = tx.Preload("Page")
	}

	incident := Incident{}
	tx = tx.Find(&incident)
	return &incident, tx.Error
}

// UpdateIncident updates a incident with the given ID.
func (c *Conn) UpdateIncident(ctx context.Context, ownerID, pageID, incidentID uint, incident *Incident) (*Incident, error) {
	if incident == nil {
		return nil, fmt.Errorf("*Incident: %w", ErrNilPointer)
	}

	i := rawIncidentWithID(ownerID, pageID, incidentID)

	tx := c.db.WithContext(ctx).Model(Incident{}).Where(&i).Updates(*incident)
	return &i, tx.Error
}

// DeleteIncident deletes the incident with the given ID.
func (c *Conn) DeleteIncident(ctx context.Context, ownerID, pageID, incidentID uint) error {
	i := rawIncidentWithID(ownerID, pageID, incidentID)

	tx := c.db.WithContext(ctx).Where(&i).Unscoped().Delete(&Incident{})
	return tx.Error
}

// checkSliceFromIDs returns a slice of raw checks from multiple IDs.
func checkSliceFromIDs(ownerID uint, checkIDs []string) []Check {
	checks := make([]Check, len(checkIDs))
	for i := range checkIDs {
		checks[i] = rawCheckWithID(ownerID, checkIDs[i])
	}

	return checks
}

// AddChecksToPage adds relationship between the checks and the page, hence
// inserting checks into the page.
func (c *Conn) AddChecksToPage(ctx context.Context, ownerID, pageID uint, checkIDs []string) error {
	if len(checkIDs) == 0 {
		return nil
	}

	p := rawPageWithID(ownerID, pageID)
	checks := checkSliceFromIDs(ownerID, checkIDs)

	return c.db.WithContext(ctx).Model(&p).Where(&p).Association("Checks").Append(checks)
}

// RemoveChecksFromPage removes relationship between the checks and the page,
// hence deleting checks from the page.
func (c *Conn) RemoveChecksFromPage(ctx context.Context, ownerID, pageID uint, checkIDs []string) error {
	if len(checkIDs) == 0 {
		return nil
	}

	p := rawPageWithID(ownerID, pageID)
	checks := checkSliceFromIDs(ownerID, checkIDs)

	return c.db.WithContext(ctx).Model(&p).Where(&p).Association("Checks").Delete(checks)
}

// rawPageTeamMemberWithID returns an empty team member with page ID and
// user ID.
func rawPageTeamMemberWithID(pageID, memberID uint, role string) PageTeam {
	pageTeamMember := PageTeam{}
	pageTeamMember.PageID = pageID
	pageTeamMember.UserID = memberID

	switch role {
	case RoleAdmin, RoleMaintainer:
		pageTeamMember.Role = role
	default:
		pageTeamMember.Role = RoleDefault
	}

	return pageTeamMember
}

// AddTeamMemberToPage adds a new team member to the page with the given ID.
func (c *Conn) AddTeamMemberToPage(ctx context.Context, ownerID, pageID, memberID uint, role string) (*PageTeam, error) {
	pt := rawPageTeamMemberWithID(pageID, memberID, role)
	p := rawPageWithID(ownerID, pageID)

	if err := c.db.WithContext(ctx).Model(&p).Where(&p).Association("Team").Append(pt); err != nil {
		return nil, err
	}

	return &pt, nil
}

// UpdateTeamMemberRole updates the role of a team member.
func (c *Conn) UpdateTeamMemberRole(ctx context.Context, ownerID, pageID, memberID uint, role string) (*PageTeam, error) {
	pt := rawPageTeamMemberWithID(pageID, memberID, role)
	p := rawPageWithID(ownerID, pageID)

	if err := c.db.WithContext(ctx).Model(&p).Where(&p).Association("Team").Replace(pt, pt); err != nil {
		return nil, err
	}

	return &pt, nil
}

// RemoveTeamMemberFromPage removes the team member from the page.
func (c *Conn) RemoveTeamMemberFromPage(ctx context.Context, ownerID, pageID, memberID uint) error {
	pt := rawPageTeamMemberWithID(pageID, memberID, "")
	p := rawPageWithID(ownerID, pageID)

	return c.db.WithContext(ctx).Model(&p).Where(&p).Association("Team").Delete(pt)
}

// GetMetricsByCheckAndStartTime fetches metrics from the metrics hypertable
// for the given check ID. It accepts a `startTime` parameter that fetches
// metrics for the check from given time.
func (c *Conn) GetMetricsByCheckAndStartTime(ctx context.Context, checkID string, startTime time.Time) ([]Metric, error) {
	metrics := []Metric{}

	tx := c.db.WithContext(ctx).Where("check_id = ? AND start_time > ?", checkID, startTime).Order("start_time DESC").Find(&metrics)
	return metrics, tx.Error
}

// GetMetricsByCheckAndDuration fetches metrics from the metrics hypertable
// for the given check ID. It accepts a `duration` parameter that fetches
// metrics for the check in the past `duration time.Duration`.
func (c *Conn) GetMetricsByCheckAndDuration(ctx context.Context, checkID string, duration time.Duration) ([]Metric, error) {
	startTime := time.Now().Add(-1 * duration)
	return c.GetMetricsByCheckAndStartTime(ctx, checkID, startTime)
}

// GetMetricsByPageAndStartTime fetches metrics for all the checks in a page
// for the given start time.
func (c *Conn) GetMetricsByPageAndStartTime(ctx context.Context, pageID uint, startTime time.Time) ([]Metric, error) {
	checkIDs := []string{}

	tx1 := c.db.WithContext(ctx).Table("page_checks").Where("page_id = ?", pageID).Pluck("check_id", &checkIDs)
	if err := tx1.Error; err != nil {
		return nil, err
	}

	metrics := []Metric{}
	tx2 := c.db.WithContext(ctx).Where("check_id IN (?) AND start_time > ?", checkIDs, startTime).
		Order("start_time DESC").
		Find(&metrics)
	return metrics, tx2.Error
}

// GetMetricsByPageAndDuration fetches metrics for all the checks in a page
// for the given duration.
func (c *Conn) GetMetricsByPageAndDuration(ctx context.Context, pageID uint, duration time.Duration) ([]Metric, error) {
	startTime := time.Now().Add(-1 * duration)
	return c.GetMetricsByPageAndStartTime(ctx, pageID, startTime)
}

// CreateMetrics inserts multiple metrics into TimescaleDB Hypertable.
func (c *Conn) CreateMetrics(ctx context.Context, metrics []Metric) error {
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
		val := fmt.Sprintf("(%s, '%s', %d, %t, %t)",
			metrics[i].CheckID,
			metrics[i].StartTime.Format(timeFormat),
			metrics[i].Duration,
			metrics[i].Timeout,
			metrics[i].Success)
		vals = append(vals, val)
	}

	args := strings.Join(vals, ", ")
	query := fmt.Sprintf(q, args)

	return c.db.WithContext(ctx).Exec(query).Error
}
