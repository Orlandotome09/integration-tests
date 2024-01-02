package individualTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator/timeTranslator"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"reflect"
	"testing"
	"time"
)

func TestTranslateName_should_return_only_first_name(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewIndividualTranslator(timetranslator)

	individual := &entity.Individual{FirstName: "AAA"}

	expected := individual.FirstName

	received := translator.TranslateName(individual)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateName_should_return_first_name_plus_last_name(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewIndividualTranslator(timetranslator)

	individual := &entity.Individual{FirstName: "AAA", LastName: "BBB"}

	expected := individual.FirstName + " " + individual.LastName

	received := translator.TranslateName(individual)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateName_should_not_translate(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewIndividualTranslator(timetranslator)

	var individual *entity.Individual

	expected := ""

	received := translator.TranslateName(individual)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateDateOfBirth_should_translate_date_of_birth_inputted(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewIndividualTranslator(timetranslator)

	dateOfBirthInputted := time.Date(2021, 9, 9, 1, 1, 1, 1, time.UTC)

	individual := &entity.Individual{DateOfBirthInputted: &dateOfBirthInputted}

	translatedDateOfBirthInputted := "02022011"

	timetranslator.On("TranslateTime", *individual.DateOfBirthInputted).Return(translatedDateOfBirthInputted)

	expected := translatedDateOfBirthInputted

	received := translator.TranslateDateOfBirth(individual)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateDateOfBirth_should_translate_date_of_birth(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewIndividualTranslator(timetranslator)

	dateOfBirth := time.Date(2021, 9, 9, 1, 1, 1, 1, time.UTC)

	individual := &entity.Individual{DateOfBirth: &dateOfBirth}

	translatedDateOfBirth := "02022011"

	timetranslator.On("TranslateTime", *individual.DateOfBirth).Return(translatedDateOfBirth)

	expected := translatedDateOfBirth

	received := translator.TranslateDateOfBirth(individual)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateDateOfBirth_should_return_empty_date_for_nil_individual(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewIndividualTranslator(timetranslator)

	var individual *entity.Individual

	expected := EmptyDateOfBirth

	received := translator.TranslateDateOfBirth(individual)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateDateOfBirth_should_return_empty_date_for_individual_without_date(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewIndividualTranslator(timetranslator)

	individual := &entity.Individual{}

	expected := EmptyDateOfBirth

	received := translator.TranslateDateOfBirth(individual)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateNationality(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewIndividualTranslator(timetranslator)

	individual := &entity.Individual{Nationality: "USA"}

	expected := individual.Nationality

	received := translator.TranslateNationality(individual)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateNationality_should_return_nationality_default_for_nil_individual(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewIndividualTranslator(timetranslator)

	var individual *entity.Individual

	expected := EmptyNationality

	received := translator.TranslateNationality(individual)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateNationality_should_return_nationality_default_for_empty_nationality(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewIndividualTranslator(timetranslator)

	individual := &entity.Individual{Nationality: ""}

	expected := EmptyNationality

	received := translator.TranslateNationality(individual)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslatePep(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewIndividualTranslator(timetranslator)

	pep := false
	individual := &entity.Individual{Pep: &pep}

	expected := *individual.Pep

	received := translator.TranslatePep(individual)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslatePep_should_return_false_for_nil_individual(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewIndividualTranslator(timetranslator)

	var individual *entity.Individual

	expected := false

	received := translator.TranslatePep(individual)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslatePep_should_return_false_for_nil_individual_without_pep_tag(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewIndividualTranslator(timetranslator)

	individual := &entity.Individual{}

	expected := false

	received := translator.TranslatePep(individual)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateIncome_should_translate_empty_income(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewIndividualTranslator(timetranslator)

	individual := &entity.Individual{}

	expected := 0.0

	received := translator.TranslateIncome(individual)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateIncome_should_translate_income(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewIndividualTranslator(timetranslator)

	income := 10.0
	individual := &entity.Individual{Income: &income}

	expected := *individual.Income

	received := translator.TranslateIncome(individual)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}

func TestTranslateIncome_should_translate_assets(t *testing.T) {
	timetranslator := &timeTranslator.MockTimeTranslator{}
	translator := NewIndividualTranslator(timetranslator)

	assets := 10.0
	individual := &entity.Individual{Assets: &assets}

	expected := *individual.Assets

	received := translator.TranslateAssets(individual)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}
