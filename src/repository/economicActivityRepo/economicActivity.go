package economicalActivityRepo

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/repository/model"
)

type economicActivityRepo struct {
	db *gorm.DB
}

func NewEconomicActivityRepo(db *gorm.DB) interfaces.EconomicActivityRepository {
	return &economicActivityRepo{
		db: db,
	}
}

func (ref *economicActivityRepo) Get(code string) (record *entity.EconomicActivity, exists bool, err error) {
	result := &model.EconomicalActivity{}

	query := ref.db.Find(result, "code_id = ?", code)
	if query.Error != nil {
		return nil, false, errors.WithStack(query.Error)
	}
	if query.RowsAffected == 0 {
		return nil, false, nil
	}

	return result.ToDomain(), true, nil
}
