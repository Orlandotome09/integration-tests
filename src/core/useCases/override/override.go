package override

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type overrideService struct {
	overrideRepository interfaces.OverrideRepository
}

func NewOverrideService(
	overrideRepository interfaces.OverrideRepository) interfaces.OverrideService {
	return &overrideService{
		overrideRepository: overrideRepository,
	}
}

func (ref *overrideService) Save(override entity.Override) error {
	if err := override.Validate(); err != nil {
		return err
	}

	overrides, err := ref.overrideRepository.FindByEntityID(override.EntityID)
	if err != nil {
		return errors.WithStack(err)
	}

	if overrides.HasInactive() {
		return values.NewErrorPrecondition(fmt.Sprintf("Entity %v is inactive.", override.EntityID))
	}

	if err := ref.overrideRepository.Save(override); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (ref *overrideService) Delete(override entity.Override) error {
	err := ref.overrideRepository.Delete(override.EntityID, override.RuleSet, override.RuleName)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (ref *overrideService) FindByEntityID(entityID uuid.UUID) (entity.Overrides, error) {
	return ref.overrideRepository.FindByEntityID(entityID)
}
