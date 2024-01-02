package values

import (
	"fmt"
)

type ProductSubCategory = string

const (
	ProductSubCategoryLegal     ProductSubCategory = "LEGAL"
	ProductSubCategoryMarketing ProductSubCategory = "MARKETING"
)

var validProductSubCategories = map[string]ProductSubCategory{
	ProductSubCategoryLegal:     ProductSubCategoryLegal,
	ProductSubCategoryMarketing: ProductSubCategoryMarketing,
}

func ParseToProductSubCategory(value string) (ProductSubCategory, error) {
	if _, in := validProductSubCategories[value]; !in {
		return "", NewErrorValidation(fmt.Sprintf("%s is an invalid product subcategory", value))
	}
	return value, nil
}
