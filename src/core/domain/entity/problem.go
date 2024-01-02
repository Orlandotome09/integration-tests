package entity

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"fmt"
	"strconv"
)

type Problem struct {
	Code   values.ProblemCode `json:"code"`
	Detail interface{}        `json:"detail"`
}

func (ref Problem) ToString() string {
	switch v := ref.Detail.(type) {
	case string:
		return ref.Detail.(string)
	case int32, int64:
		return strconv.Itoa(ref.Detail.(int))
	case float32, float64:
		return strconv.FormatFloat(ref.Detail.(float64), 'E', -1, 64)
	case map[string]interface{}, []string:
		str, _ := json.Marshal(ref.Detail)
		return string(str)
	default:
		fmt.Sprintf("%v", v)
		return ref.Code
	}
}
