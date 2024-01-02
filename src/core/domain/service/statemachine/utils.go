package statemachine

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type profileStateBuilder struct {
	state *entity.State
}

type overrideBuilder struct {
	overrides entity.Overrides
}

func newProfileStateBuilder() profileStateBuilder {
	return profileStateBuilder{
		state: &entity.State{},
	}
}

func (ref profileStateBuilder) withResult(result values.Result, pending bool, skipForApproval bool) profileStateBuilder {
	ref.state.ValidationStepsResults = []entity.ValidationStepResult{
		{
			Result:          "",
			StepNumber:      0,
			SkipForApproval: skipForApproval,
			RuleResults: []entity.RuleResultV2{
				{
					RuleSet:  values.RuleSetTest,
					RuleName: values.RuleNameTest,
					Result:   result,
					ExpireAt: nil,
					Metadata: nil,
					Pending:  pending,
					Tags:     nil,
					Problems: nil,
				},
			},
		},
	}

	return ref
}

func (ref profileStateBuilder) withAdditionalStep(result values.Result, pending bool, skipForApproval bool) profileStateBuilder {
	ref.state.ValidationStepsResults = append(ref.state.ValidationStepsResults,
		entity.ValidationStepResult{
			Result:          "",
			StepNumber:      1,
			SkipForApproval: skipForApproval,
			RuleResults: []entity.RuleResultV2{
				{
					RuleSet:  values.RuleSetTest,
					RuleName: values.RuleNameTest,
					Result:   result,
					ExpireAt: nil,
					Metadata: nil,
					Pending:  pending,
					Tags:     nil,
					Problems: nil,
				},
			},
		})

	return ref
}

func newOverrideBuilder() overrideBuilder {
	return overrideBuilder{
		overrides: entity.Overrides{},
	}
}

func (ref overrideBuilder) withOverride(result values.Result) overrideBuilder {
	ref.overrides = append(ref.overrides, entity.Override{
		RuleSet:  values.RuleSetTest,
		RuleName: values.RuleNameTest,
		Result:   result,
	})

	return ref
}

func (ref overrideBuilder) isBlocked() overrideBuilder {
	ref.overrides = append(ref.overrides, entity.Override{
		RuleSet:  values.RuleSetState,
		RuleName: values.RuleNameBlocked,
		Result:   values.ResultStatusRejected,
	})

	return ref
}

func (ref overrideBuilder) isInactive() overrideBuilder {
	ref.overrides = append(ref.overrides, entity.Override{
		RuleSet:  values.RuleSetState,
		RuleName: values.RuleNameInactive,
		Result:   values.ResultStatusRejected,
	})

	return ref
}

func (ref profileStateBuilder) build() *entity.State {
	return ref.state
}

func (ref overrideBuilder) build() entity.Overrides {
	return ref.overrides
}
