package infra

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenConnection(host, port, user, dbname, password string, sslMode bool) (*gorm.DB, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", host, port, user, dbname, password)

	if !sslMode {
		connString = connString + " sslmode=disable"
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: connString,
	}), &gorm.Config{
		PrepareStmt: false,
	})
	if err != nil {
		return nil, errors.WithMessagef(err,
			"Failed to connect database, host: %s - port: %s - user: %s - dbname: %s",
			host, port, user, dbname)
	}

	sqlDB, _ := db.DB()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(50)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	logrus.Infof("Connected to database, host: %s - port: %s - user: %s - dbname: %s", host, port, user, dbname)

	return db, nil
}
