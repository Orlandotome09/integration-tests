package override

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	overrideRepository *mocks.OverrideRepository
	service            interfaces.OverrideService
)

func TestMain(m *testing.M) {
	overrideRepository = &mocks.OverrideRepository{}
	service = NewOverrideService(overrideRepository)
	os.Exit(m.Run())
}

func TestSave_Active(t *testing.T) {
	parentID := uuid.New()

	override := entity.Override{
		EntityID:   uuid.New(),
		EntityType: values.EntityTypeProfile,
		ParentID:   &parentID,
		RuleSet:    values.RuleSetBlacklist,
		RuleName:   values.RuleNameBlocked,
		Result:     values.ResultStatusAnalysing,
	}
	overrides := entity.Overrides{}

	overrideRepository.On("FindByEntityID", override.EntityID).Return(overrides, nil)
	overrideRepository.On("Save", override).Return(nil)

	err := service.Save(override)

	assert.Nil(t, err)
}

func TestSave_Inactive(t *testing.T) {
	parentID := uuid.New()

	override := entity.Override{
		EntityID:   uuid.New(),
		EntityType: values.EntityTypeProfile,
		ParentID:   &parentID,
		RuleSet:    values.RuleSetBlacklist,
		RuleName:   values.RuleNameDocumentNotFound,
		Result:     values.ResultStatusCreated,
	}
	overrides := entity.Overrides{}
	overrides = append(overrides, entity.Override{EntityID: uuid.New(), EntityType: values.EntityTypeProfile, RuleSet: values.RuleSetState,
		RuleName: values.RuleNameInactive, Result: values.ResultStatusInactive})

	overrideRepository.On("FindByEntityID", override.EntityID).Return(overrides, nil)

	expected := values.NewErrorPrecondition(fmt.Sprintf("Entity %v is inactive.", override.EntityID))
	err := service.Save(override)

	assert.Equal(t, expected, err)
}

func TestDelete(t *testing.T) {
	parentID := uuid.New()

	override := entity.Override{
		EntityID:   uuid.New(),
		EntityType: values.EntityTypeProfile,
		ParentID:   &parentID,
		RuleSet:    values.RuleSetState,
		RuleName:   values.RuleNameInactive,
		Result:     values.ResultStatusInactive,
	}

	overrideRepository.On("Delete", override.EntityID, override.RuleSet, override.RuleName).Return(nil)

	err := service.Delete(override)

	assert.Nil(t, err)
}

func TestFindByEntityID(t *testing.T) {
	entityID := uuid.New()
	overrides := entity.Overrides{entity.Override{}, entity.Override{}}

	overrideRepository.On("FindByEntityID", entityID).Return(overrides, nil)

	expected := entity.Overrides{entity.Override{}, entity.Override{}}
	received, err := service.FindByEntityID(entityID)

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}
