package engines

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type personSubEngine struct {
	person           entity.Person
	ruleValidator    interfaces.RuleValidator
	personFactory    interfaces.PersonFactory
	personRepository interfaces.PersonRepository
}

func NewPersonSubEngine(ruleValidator interfaces.RuleValidator,
	personFactory interfaces.PersonFactory,
	personRepository interfaces.PersonRepository) interfaces.SubEngine {
	return &personSubEngine{
		ruleValidator:    ruleValidator,
		personFactory:    personFactory,
		personRepository: personRepository,
	}
}

func (ref *personSubEngine) Prepare(person entity.Person, offerType string) error {
	ref.person = person

	result, err := ref.personFactory.Build(person)
	if err != nil {
		return errors.WithStack(err)
	}

	if result != nil {
		ref.person = *result
	}

	if ref.person.CadastralValidationConfig == nil {
		return errors.New(fmt.Sprintf("[personSubEngine] Catalog not found. Person id: %v", person.EntityID.String()))
	}

	ref.ruleValidator.SetRules(ref.person.ValidationSteps)

	_, err = ref.personRepository.Save(ref.person)
	if err != nil {
		return errors.WithStack(err)
	}

	logrus.WithField("person", ref.person).Info("[personSubEngine] Person built")

	return nil
}

func (ref *personSubEngine) Validate(state entity.State, override entity.Overrides, noCache bool, entityID uuid.UUID, engineName string) (*entity.State, error) {

	return ref.ruleValidator.Validate(state, override, noCache, entityID, engineName)
}

func (ref *personSubEngine) PosProcessing(previousState *entity.State, newState *entity.State, entityID uuid.UUID) error {
	return nil
}

func (ref *personSubEngine) GetName() string {
	return values.EngineNamePerson
}

func (ref *personSubEngine) NewInstance() interfaces.SubEngine {
	newPersonSubEngine := &personSubEngine{
		ruleValidator:    ref.ruleValidator.NewInstance(),
		personFactory:    ref.personFactory,
		personRepository: ref.personRepository,
	}

	return newPersonSubEngine
}
