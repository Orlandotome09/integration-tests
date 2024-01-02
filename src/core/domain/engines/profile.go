package engines

import (
	"fmt"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type profileEngine struct {
	profile           entity.Profile
	validator         interfaces.RuleValidator
	posProcessor      interfaces.PosProcessor
	profileAdapter    interfaces.ProfileAdapter
	profileFactory    interfaces.ProfileFactory
	profileRepository interfaces.ComplianceProfileRepository
}

func NewProfileEngine(validator interfaces.RuleValidator,
	posProcessor interfaces.PosProcessor,
	profileService interfaces.ProfileAdapter,
	factory interfaces.ProfileFactory,
	profileRepository interfaces.ComplianceProfileRepository,
) interfaces.Engine {
	return &profileEngine{
		validator:         validator,
		posProcessor:      posProcessor,
		profileAdapter:    profileService,
		profileFactory:    factory,
		profileRepository: profileRepository,
	}
}

func (ref *profileEngine) Prepare(entityID uuid.UUID) error {
	profile, err := ref.profileFactory.Build(entityID)
	if err != nil {
		return errors.WithStack(err)
	}

	if profile == nil {
		return values.NewErrorNotFound("profile")
	}

	ref.profile = *profile

	if !profile.HasCadastralValidationConfig() {
		return errors.New(fmt.Sprintf("[profileEngine] cadastral validation config not found for profile id: %v", entityID.String()))
	}

	ref.validator.SetRules(profile.ValidationSteps)

	_, err = ref.profileRepository.Save(*profile)
	if err != nil {
		return errors.WithStack(err)
	}

	logrus.WithField("profile", ref.profile).Info("[profileEngine] Profile built")

	return nil
}

func (ref *profileEngine) Validate(state entity.State, override entity.Overrides, noCache bool, entityID uuid.UUID, engineName string) (*entity.State, error) {

	return ref.validator.Validate(state, override, noCache, entityID, engineName)
}

func (ref *profileEngine) PosProcessing(previousState *entity.State, newState *entity.State, entityID uuid.UUID) error {

	if newState.IsApproved() {
		if !ref.profile.HasProductConfig() {
			logrus.
				WithField("profile", ref.profile).
				Infof("[profile.PosProcessing] Profile, Catalog or Product Config are not set for profile. Skipping post processing. ProfileID: %v", newState.EntityID)
			return nil
		}

		if ref.profile.ShouldCreateInternalAccount() {
			if err := ref.posProcessor.CreateInternalAccount(entityID); err != nil {
				return errors.WithStack(err)
			}
		}

		if err := ref.posProcessor.SendToLimit(ref.profile, *newState, *ref.profile.CadastralValidationConfig.ProductConfig); err != nil {
			return errors.WithStack(err)
		}

		if err := ref.posProcessor.SendToTreeAdapter(ref.profile, *ref.profile.CadastralValidationConfig.ProductConfig); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (ref *profileEngine) GetName() string {
	return values.EngineNameProfile
}

func (ref *profileEngine) NewInstance() interfaces.Engine {
	newProfileEngine := &profileEngine{
		validator:         ref.validator.NewInstance(),
		posProcessor:      ref.posProcessor,
		profileFactory:    ref.profileFactory,
		profileAdapter:    ref.profileAdapter,
		profileRepository: ref.profileRepository,
	}
	return newProfileEngine
}
