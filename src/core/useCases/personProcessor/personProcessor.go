package personProcessor

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

type compliancePersonProcessor struct {
	personSubEngine    interfaces.SubEngine
	complianceAnalyzer interfaces.ComplianceAnalyzer
}

func NewCompliancePersonProcessor(
	personSubEngine interfaces.SubEngine,
	complianceAnalyzer interfaces.ComplianceAnalyzer) interfaces.CompliancePersonProcessor {
	return &compliancePersonProcessor{
		personSubEngine:    personSubEngine,
		complianceAnalyzer: complianceAnalyzer,
	}
}

func (ref *compliancePersonProcessor) ExecuteForPerson(person entity.Person, offer string) (*entity.State, error) {

	subEngineName := person.EntityType
	now := time.Now()

	logrus.
		WithField("engine_name", subEngineName).
		WithField("entity_id", person.EntityID).
		WithField("date", now).
		Info("Running Compliance Analysis for Person")

	newSubEngine := ref.personSubEngine.NewInstance()

	err := newSubEngine.Prepare(person, offer)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	state, err := ref.complianceAnalyzer.RunComplianceAnalysis(newSubEngine, person.EntityID, now, now, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return state, nil
}
