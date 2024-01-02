package timeTranslator

import (
	"fmt"
	"strconv"
	"time"
)

type TimeTranslator interface {
	TranslateTime(dateTime time.Time) string
	GenerateMessageDate() time.Time
}

type timeTranslator struct{}

func NewTimeTranslator() TimeTranslator {
	return &timeTranslator{}
}

func (ref *timeTranslator) TranslateTime(dateTime time.Time) string {
	day := dateTime.Day()
	month := int(dateTime.Month())
	year := dateTime.Year()

	dayStr, monthStr, yearStr := strconv.Itoa(day), strconv.Itoa(month), strconv.Itoa(year)
	if day/10 == 0 {
		dayStr = "0" + strconv.Itoa(day)
	}

	if month/10 == 0 {
		monthStr = "0" + strconv.Itoa(month)
	}

	return fmt.Sprintf("%s%s%s", dayStr, monthStr, yearStr)
}

func (ref *timeTranslator) GenerateMessageDate() time.Time {
	return time.Now().UTC()
}
