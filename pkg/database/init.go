// Package database contains the connection to database for the status app
// as well as the timeseries db (timescale) for storing the metrics. It
// contains methods and types to interact with the database using an ORM.
package database

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // PostgreSQL
	log "github.com/sirupsen/logrus"

	"github.com/sdslabs/status/pkg/config"
	"github.com/sdslabs/status/pkg/controller"
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

func setupMetricsExporter(conf *metrics.ProviderConfig, manager *controller.Manager) (*metrics.TimescaleExporter, error) {
	if err := initFromProvider(conf); err != nil {
		return nil, err
	}
	uid, err := createStandaloneUser()
	if err != nil {
		return nil, err
	}
	return &metrics.TimescaleExporter{
		Manager:  manager,
		UserID:   uid,
		Interval: conf.Interval,
		Quit:     make(chan bool),
	}, nil
}

func setupMetrics(ex *metrics.TimescaleExporter) {
	ticker := time.NewTicker(ex.Interval)
	for {
		select {
		case <-ticker.C:
			log.Infoln("Trying to insert metrics into DB")
			stats := ex.PullLatestControllerStatistics()
			metricsToInsert := []Metric{}
			for _, st := range stats {
				checkID, err := strconv.Atoi(st.Name)
				if err != nil {
					log.Errorf("Failed to insert metric for check=%s", st.Name)
					continue
				}
				metricsToInsert = append(metricsToInsert, Metric{
					CheckID:   uint(checkID),
					StartTime: &st.StartTime,
					Duration:  st.Duration,
					Timeout:   st.Timeout,
					Success:   st.Success,
				})
			}
			if err := CreateMetrics(metricsToInsert); err != nil {
				log.Errorf("Error while inserting metrics to DB: %s", err.Error())
				ex.ErrCount++
			}
		case <-ex.Quit:
			log.Infoln("Stopping timescale metrics exporter")
			close(ex.Quit)
			return
		default:
			if ex.ErrCount >= 10 {
				log.Fatalf("Stopping metrics collector, too many errors!")
			}
			continue
		}
	}
}

// SetupConf defines the setup configuration for timescale metrics.
type SetupConf struct {
	*metrics.ProviderConfig
	*controller.Manager
	Standalone bool
	Checks     []*config.CheckConf
}

// SetupMetrics sets up the timescale metrics for agent.
func SetupMetrics(conf *SetupConf) {
	ex, err := setupMetricsExporter(conf.ProviderConfig, conf.Manager)
	if err != nil {
		log.Fatalf("Cannot setup metrics: %s", err.Error())
		return
	}

	// Insert checks into DB before starting the metrics collector for standalone mode.
	if conf.Standalone {
		for _, checkConfig := range conf.Checks {
			payloads := []Payload{}
			for _, payload := range checkConfig.GetPayloads() {
				payloads = append(payloads, Payload{
					Type:  payload.GetType(),
					Value: payload.GetValue(),
				})
			}
			check := &Check{
				OwnerID:     ex.UserID,
				Interval:    int(checkConfig.GetInterval()),
				Timeout:     int(checkConfig.GetTimeout()),
				InputType:   checkConfig.GetInput().GetType(),
				InputValue:  checkConfig.GetInput().GetValue(),
				OutputType:  checkConfig.GetOutput().GetType(),
				OutputValue: checkConfig.GetOutput().GetValue(),
				TargetType:  checkConfig.GetTarget().GetType(),
				TargetValue: checkConfig.GetTarget().GetValue(),
				Title:       checkConfig.GetName(),
				Payloads:    payloads,
			}
			log.Infof("Inserting check '%s' inside DB", check.Title)
			createdCheck, err := CreateCheck(check)
			if err != nil {
				log.Errorf("Failed to insert check '%s'", check.Title)
				log.Errorln("Skipping and moving ahead...")
				continue
			}
			checkConfig.ID = createdCheck.ID
		}
	}

	go setupMetrics(ex)
}
