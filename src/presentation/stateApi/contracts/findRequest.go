package contracts

import (
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type FindRequest struct {
	ProfileID    string             `uri:"profile_id" binding:"uuid"`
	ResultStatus values2.Result     `form:"result_status"`
	EngineName   values2.EngineName `form:"engine_name"`
}
