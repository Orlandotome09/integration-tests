package repository

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"bitbucket.org/bexstech/temis-compliance/src/repository/model"
)

func Migrate(db *gorm.DB) error {
	err := prevMigrations(db)
	if err != nil {
		return errors.WithStack(err)
	}

	err = fixedMigrations(db)
	if err != nil {
		return errors.WithStack(err)
	}

	err = temporaryMigrations(db)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func prevMigrations(db *gorm.DB) error {

	return nil
}

func fixedMigrations(db *gorm.DB) error {
	autoMigrateResult := db.AutoMigrate(
		&model.DOAResult{},
		&model.Override{},
		&model.State{},
		&model.Offer{},
		&model.EconomicalActivity{},
		&model.Contract{},
		&model.ComplianceProfile{},
		&model.Person{},
	)
	if autoMigrateResult != nil {
		return errors.WithStack(autoMigrateResult)
	}

	migrateResult := db.Model(&model.Override{})
	if migrateResult.Error != nil {
		return errors.WithStack(migrateResult.Error)
	}

	migrateResult = db.Model(&model.DOAResult{})
	if migrateResult.Error != nil {
		return errors.WithStack(migrateResult.Error)
	}

	return nil
}

func temporaryMigrations(db *gorm.DB) error {
	return nil
}
