package doaResultRepo

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/repository/model"
)

type doaRepository struct {
	_interfaces.DOAResultRepository
	db *gorm.DB
}

func NewDOAResultRepository(db *gorm.DB) _interfaces.DOAResultRepository {
	return &doaRepository{db: db}
}

func (ref *doaRepository) Get(id *uuid.UUID) (*entity.DOAResult, error) {
	record, err := ref.get(id)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if record == nil {
		return nil, nil
	}

	result := record.ToDomain()
	return &result, nil
}

func (ref *doaRepository) get(id *uuid.UUID) (*model.DOAResult, error) {
	record := &model.DOAResult{}

	query := ref.db.Find(record, "id = ?", *id)
	if query.Error != nil {
		return nil, errors.WithStack(query.Error)
	}
	if query.RowsAffected == 0 {
		return nil, nil
	}

	return record, nil
}

func (ref *doaRepository) Save(doaResult *entity.DOAResult) (*entity.DOAResult, error) {
	record := &model.DOAResult{}
	record.FromDomain(doaResult)

	if db := ref.db.Save(&record); db.Error != nil {
		return nil, errors.WithStack(db.Error)
	}

	result := record.ToDomain()
	return &result, nil
}

func (ref *doaRepository) Enrich(doaResult *entity.DOAResult) (*entity.DOAResult, error) {
	record := &model.DOAResult{}
	record.FromDomain(doaResult)

	current, err := ref.get(&record.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	record.PrepareWithCurrent(current)

	if db := ref.db.Save(&record); db.Error != nil {
		return nil, errors.WithStack(db.Error)
	}
	result := record.ToDomain()
	return &result, nil
}

func (ref *doaRepository) FindByEntityID(entityID *uuid.UUID) ([]entity.DOAResult, error) {
	var records []model.DOAResult

	query := ref.db.Find(&records, "entity_id = ?", *entityID)
	if query.Error != nil {
		return nil, errors.WithStack(query.Error)
	}
	if query.RowsAffected == 0 {
		return nil, nil
	}

	doaResults := make([]entity.DOAResult, len(records))
	for index, record := range records {
		doaResults[index] = record.ToDomain()
	}
	return doaResults, nil
}

func (ref *doaRepository) FindByEntityIdAndDocumentId(entityID *uuid.UUID,
	documentID *uuid.UUID) ([]entity.DOAResult, error) {
	var records []model.DOAResult

	query := ref.db.Find(&records, "entity_id = ? and document_id = ?", *entityID, *documentID)
	if query.Error != nil {
		return nil, errors.WithStack(query.Error)
	}
	if query.RowsAffected == 0 {
		return nil, nil
	}

	doaResults := []entity.DOAResult{}
	for _, record := range records {
		doaResult := record.ToDomain()
		doaResults = append(doaResults, doaResult)
	}

	return doaResults, nil
}

func (ref *doaRepository) FindLastByEntityIdAndDocumentId(entityID *uuid.UUID,
	documentID *uuid.UUID) (*entity.DOAResult, error) {
	var records []model.DOAResult

	query := ref.db.Order("created_at DESC").Limit(1).Find(&records, "entity_id = ? and document_id = ?", *entityID, *documentID)
	if query.Error != nil {
		return nil, errors.WithStack(query.Error)
	}
	if query.RowsAffected == 0 {
		return nil, nil
	}

	result := records[0].ToDomain()

	return &result, nil
}
