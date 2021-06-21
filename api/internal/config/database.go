package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
	"github.com/pkg/errors"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
)

const (
	retryDBAfter = 5  // seconds
	failDBAfter  = 30 // seconds
	dbServer     = "catalogue"
)

// GetDBConnection returns postgress db connection
func (config *Config) GetDBConnection(connStr string) (*sql.DB, error) {
	if config.DatadogEnableAPM {
		sqltrace.Register(dbServer, &pq.Driver{}, sqltrace.WithServiceName(config.DatadogServiceName))
	} else {
		sql.Register(dbServer, &pq.Driver{})
	}
	db, err := config.tryPG(connStr, retryDBAfter, failDBAfter)
	if err != nil {
		log.Fatal(err)
	}
	return db, nil

}

// MigrateFolder migrates db
func (config *Config) MigrateFolder(connStr string, path string) {
	db, err := config.tryPG(connStr, retryDBAfter, failDBAfter)
	if err != nil {
		log.Fatalln("MigrateFolder", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed to instantiate db migration driver"))
	}
	mig, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", path), "postgres", driver)
	if err != nil {
		log.Fatalln(errors.Wrap(err, fmt.Sprintf("failed to instantiate migration instance: %s", path)))
	}
	if err := mig.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalln(errors.Wrap(err, fmt.Sprintf("error while running initial migration up to setup db %s", path)))
	}
	err1, err2 := mig.Close()
	if err1 != nil {
		log.Fatalln(errors.Wrap(err1, fmt.Sprintf("could not close migrate source: %s", path)))
	}
	if err2 != nil {
		log.Fatalln(errors.Wrap(err2, fmt.Sprintf("could not close migrate database: %s", path)))
	}
	log.Println("migrated")
	return
}

// GetDBConnectionString returns the database connection string based on config values
func (config *Config) GetDBConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s%s/%s?sslmode=disable",
		config.DBUserName,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)
}

// ConnectPG to tries setting up a db pool, and pinging the db to be sure it's
// running correctly.
func (config *Config) ConnectPG(connString string) (*sql.DB, error) {
	var db *sql.DB
	var err error
	if config.DatadogEnableAPM {
		db, err = sqltrace.Open(dbServer, connString, sqltrace.WithServiceName(config.DatadogServiceName))
	} else {
		db, err = sql.Open(dbServer, connString)
	}
	// db, err := sql.Open("postgres", connString)
	// if err != nil {
	// 	return nil, err
	// }
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrapf(err,
			"failed to connect to %s",
		)
	}
	return db, nil
}

// TryPG tries to connect to a Postgres database using the give connection
// string. If the connection fails, it will continue trying every 'wait' seconds
// until 'max' is reached.
func (config *Config) tryPG(connString string, wait, max int) (*sql.DB, error) {
	db, err := config.ConnectPG(connString)
	if err == nil {
		return db, nil
	}

	ticker := time.NewTicker(time.Duration(wait) * time.Second)
	defer ticker.Stop()
	done := make(chan bool)
	go func() {
		time.Sleep(time.Duration(max) * time.Second)
		done <- true
	}()
	for {
		select {
		case <-done:
			return nil, errors.Wrapf(err,
				"failed to connect to Postgres within %d seconds",
				max,
			)
		case <-ticker.C:
			log.Println("retrying database connection..")
			db, err = config.ConnectPG(connString)
			if err == nil {
				return db, nil
			}
		}
	}
}
