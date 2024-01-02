package contractRepo

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/repository/errorHandler"
	"bitbucket.org/bexstech/temis-compliance/src/repository/model"
)

type contractRepository struct {
	db           *gorm.DB
	errorHandler errorHandler.ErrorHandler
}

func New(db *gorm.DB, errorHandler errorHandler.ErrorHandler) interfaces.ContractRepository {
	return &contractRepository{
		db:           db,
		errorHandler: errorHandler,
	}
}

func (ref *contractRepository) Get(id uuid.UUID) (*entity.Contract, error) {
	record := model.Contract{}
	result := ref.db.Find(&record, "contract_id = ?", id)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	contract := record.ToDomain()
	return &contract, nil
}

func (ref *contractRepository) Save(contract entity.Contract) (*entity.Contract, error) {
	record := model.NewContractFromDomain(contract)
	if err := ref.db.Save(&record).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return &contract, nil
}
