package companyTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator/timeTranslator"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"reflect"
	"testing"
	"time"
)

func TestTranslateLegalName(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewCompanyTranslator(timetranslator)

	company := &entity.Company{LegalName: "AAA"}

	expected := "AAA"

	received := translator.TranslateLegalName(company)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateLegalName_should_return_empty_for_nil_company(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewCompanyTranslator(timetranslator)

	var company *entity.Company

	expected := ""

	received := translator.TranslateLegalName(company)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateDateOfIncorporation(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewCompanyTranslator(timetranslator)

	dateOfIncorporation := time.Date(2021, 9, 9, 1, 1, 1, 1, time.UTC)

	company := &entity.Company{LegalName: "AAA", DateOfIncorporation: &dateOfIncorporation}

	translatedDateOfIncorporation := "02022011"

	timetranslator.On("TranslateTime", *company.DateOfIncorporation).Return(translatedDateOfIncorporation)

	expected := translatedDateOfIncorporation

	received := translator.TranslateDateOfIncorporation(company)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateDateOfIncorporation_should_return_default_for_nil_company(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewCompanyTranslator(timetranslator)

	var company *entity.Company

	expected := EmptyTranslatedDateOfIncorporation

	received := translator.TranslateDateOfIncorporation(company)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateDateOfIncorporation_should_return_empty_for_company_without_date(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewCompanyTranslator(timetranslator)

	company := &entity.Company{}

	expected := EmptyTranslatedDateOfIncorporation

	received := translator.TranslateDateOfIncorporation(company)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslatePlaceOfIncorporation(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewCompanyTranslator(timetranslator)

	company := &entity.Company{PlaceOfIncorporation: "USA"}

	expected := company.PlaceOfIncorporation

	received := translator.TranslatePlaceOfIncorporation(company)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslatePlaceOfIncorporation_should_return_default_for_nil_company(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewCompanyTranslator(timetranslator)

	var company *entity.Company

	expected := EmptyTranslatedPlaceOfIncorporation

	received := translator.TranslatePlaceOfIncorporation(company)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslatePlaceOfIncorporation_should_return_default_for_company_without_place_of_incorporation(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewCompanyTranslator(timetranslator)

	company := &entity.Company{}

	expected := EmptyTranslatedPlaceOfIncorporation

	received := translator.TranslatePlaceOfIncorporation(company)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateIncome_should_not_translate_for_company_nil(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewCompanyTranslator(timetranslator)

	var company *entity.Company

	expected := 0.0

	received := translator.TranslateIncome(company)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateIncome_should_translate_income(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewCompanyTranslator(timetranslator)

	company := &entity.Company{AnnualIncome: 120.0}

	expected := company.AnnualIncome / 12

	received := translator.TranslateIncome(company)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}
