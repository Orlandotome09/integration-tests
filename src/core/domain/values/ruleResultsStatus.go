package values

import (
	"fmt"
)

type Result string

const (
	ResultStatusCreated    Result = "CREATED"
	ResultStatusIncomplete Result = "INCOMPLETE"
	ResultStatusAnalysing  Result = "ANALYSING"
	ResultStatusApproved   Result = "APPROVED"
	ResultStatusRejected   Result = "REJECTED"
	ResultStatusIgnored    Result = "IGNORED"
	ResultStatusBlocked    Result = "BLOCKED"
	ResultStatusInactive   Result = "INACTIVE"
)

var validResultStatus = map[string]Result{
	ResultStatusCreated.ToString():    ResultStatusCreated,
	ResultStatusIncomplete.ToString(): ResultStatusIncomplete,
	ResultStatusAnalysing.ToString():  ResultStatusAnalysing,
	ResultStatusApproved.ToString():   ResultStatusApproved,
	ResultStatusRejected.ToString():   ResultStatusRejected,
	ResultStatusIgnored.ToString():    ResultStatusIgnored,
	ResultStatusBlocked.ToString():    ResultStatusBlocked,
	ResultStatusInactive.ToString():   ResultStatusInactive,
}

func (result Result) Validate() error {
	_, in := validResultStatus[result.ToString()]
	if !in {
		return NewErrorValidation(fmt.Sprintf("%s is an invalid result status", result))
	}
	return nil
}

func (result Result) ToString() string {
	return string(result)
}

func (result *Result) IsItWorstThan(otherResult *Result) bool {

	return bestToWorstResultStatusOrder[*result] > bestToWorstResultStatusOrder[*otherResult]
}

var bestToWorstResultStatusOrder = map[Result]int{
	ResultStatusApproved:   1,
	ResultStatusAnalysing:  2,
	ResultStatusRejected:   3,
	ResultStatusIncomplete: 4,
	ResultStatusBlocked:    5,
	ResultStatusInactive:   6,
}
