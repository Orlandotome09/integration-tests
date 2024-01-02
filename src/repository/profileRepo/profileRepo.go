package profileRepo

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/repository/errorHandler"
	"bitbucket.org/bexstech/temis-compliance/src/repository/model"
)

type profileRepository struct {
	db           *gorm.DB
	errorHandler errorHandler.ErrorHandler
}

func New(db *gorm.DB, errorHandler errorHandler.ErrorHandler) interfaces.ComplianceProfileRepository {
	return &profileRepository{
		db:           db,
		errorHandler: errorHandler,
	}
}

func (ref *profileRepository) Save(profile entity.Profile) (*entity.Profile, error) {
	record, err := model.NewProfileFromDomain(profile)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if result := ref.db.Save(record); result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	return ref.Get(record.ProfileID)

}

func (ref *profileRepository) Get(profileID uuid.UUID) (*entity.Profile, error) {
	record := model.ComplianceProfile{}

	result := ref.db.Find(&record, "profile_id = ?", profileID)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	profile, err := record.ToDomain()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return profile, nil
}
