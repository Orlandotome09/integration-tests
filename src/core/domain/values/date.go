package values

import (
	"encoding/json"
	"fmt"
	"time"
)

type Date time.Time

var _ json.Unmarshaler = &Date{}
var _ json.Marshaler = &Date{}

func (ref *Date) UnmarshalJSON(bs []byte) error {
	var s string
	err := json.Unmarshal(bs, &s)
	if err != nil {
		return err
	}
	t, err := time.ParseInLocation("2006-01-02", s, time.UTC)
	if err != nil {
		return err
	}
	*ref = Date(t)
	return nil
}

func (ref Date) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(ref).Format("2006-01-02"))
	return []byte(stamp), nil
}

func (ref *Date) ToTime() *time.Time {
	if ref == nil {
		return nil
	}
	newTime := time.Time(*ref)
	return &newTime
}

func NewDateFromTime(tim *time.Time) *Date {
	if tim == nil {
		return nil
	}

	date := Date(*tim)
	return &date
}
