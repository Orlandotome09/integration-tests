package overrideRepo

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	values "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bitbucket.org/bexstech/temis-compliance/src/repository/model"
)

type overrideSqlRepository struct {
	db *gorm.DB
}

func NewOverrideSqlRepository(db *gorm.DB) interfaces.OverrideRepository {
	return &overrideSqlRepository{db: db}
}

func (ref *overrideSqlRepository) Save(override entity.Override) error {
	record := model.Override{}.FromDomain(override)

	if result := ref.db.Save(&record); result.Error != nil {
		return result.Error
	}

	return nil
}

func (ref *overrideSqlRepository) Delete(entityID uuid.UUID, ruleSet values.RuleSet,
	ruleName values.RuleName) error {
	input := model.Override{
		EntityID: entityID,
		RuleSet:  ruleSet.ToString(),
		RuleName: ruleName.ToString(),
	}

	result := ref.db.Unscoped().Delete(&input)
	if result.Error != nil {
		return errors.WithStack(result.Error)
	}

	if result.RowsAffected == 0 {
		return values.NewErrorNotFound(
			fmt.Sprintf("Override for key (%v) (%v) (%v)", input.EntityID, input.RuleSet, input.RuleName))
	}

	return nil
}

func (ref *overrideSqlRepository) FindByEntityID(entityID uuid.UUID) (entity.Overrides, error) {
	var records []model.Override

	if result := ref.db.Where("entity_id = ?", entityID).Find(&records); result.Error != nil {
		return entity.Overrides{}, errors.WithStack(result.Error)
	}

	var overrides entity.Overrides

	for _, record := range records {
		override := record.ToDomain()
		overrides = append(overrides, override)
	}

	return overrides, nil
}
