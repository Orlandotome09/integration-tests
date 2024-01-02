package values

import (
	"fmt"
)

type FileSide string

const (
	FileSideFront     FileSide = "FRONT"
	FileSideBack      FileSide = "BACK"
	FileSideBackFront FileSide = "ALL"
)

var validFileSides = map[string]FileSide{
	FileSideFront.ToString():     FileSideFront,
	FileSideBack.ToString():      FileSideBack,
	FileSideBackFront.ToString(): FileSideBackFront,
}

func (fileSide FileSide) Validate() error {
	_, in := validFileSides[fileSide.ToString()]
	if !in {
		return NewErrorValidation(fmt.Sprintf("%s is an invalid file side", fileSide))
	}
	return nil
}

func (fileSide FileSide) ToString() string {
	return string(fileSide)
}
