package personRepo

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/repository/errorHandler"
	"bitbucket.org/bexstech/temis-compliance/src/repository/model"
)

type personRepository struct {
	db           *gorm.DB
	errorHandler errorHandler.ErrorHandler
}

func New(db *gorm.DB, errorHandler errorHandler.ErrorHandler) interfaces.PersonRepository {
	return &personRepository{
		db:           db,
		errorHandler: errorHandler,
	}
}

func (ref *personRepository) Save(person entity.Person) (*entity.Person, error) {
	record, err := model.NewPersonFromDomain(person)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if result := ref.db.Save(record); result.Error != nil {
		if ref.errorHandler.IsRecordDuplicated(result) {
			return ref.Get(record.EntityID)
		}
		return nil, errors.WithStack(result.Error)
	}

	return ref.Get(record.EntityID)
}

func (ref *personRepository) Get(personID uuid.UUID) (*entity.Person, error) {
	record := model.Person{}

	result := ref.db.Find(&record, "entity_id = ?", personID)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	person, err := record.ToDomain()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return person, nil
}
