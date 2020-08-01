package service_test

import (
	"strings"
	"testing"

	"github.com/Bocmah/phpdocker-scaffold/pkg/service"
)

type validationResult struct {
	wantErrs     []string
	actualErrs   error
	validatedVal interface{}
}

func failTestOnUnspottedError(result validationResult, t *testing.T) {
	for _, e := range result.wantErrs {
		if !strings.Contains(result.actualErrs.Error(), e) {
			t.Errorf("Failed to spot error %s in value %v", e, result.validatedVal)

		}
	}
}

func failTestOnErrorsOnCorrectInput(errs error, t *testing.T) {
	if errs != nil {
		t.Errorf("Following errors were returned despite correct inputs %v", errs)
	}
}

func TestDatabase_FillDefaultsIfNotSet(t *testing.T) {
	db := service.DatabaseConfig{}

	db.FillDefaultsIfNotSet()

	want := service.DatabaseConfig{
		System:  service.MySQL,
		Port:    3306,
		Version: "8.0",
	}

	if db != want {
		t.Errorf("Incorrect defaults, want %v, got %v", want, db)
	}
}

func TestDatabase_ValidateIncorrectInput(t *testing.T) {
	db := service.DatabaseConfig{System: "Unsupported"}

	errs := db.Validate()

	if errs != nil {
		res := validationResult{
			wantErrs: []string{
				"Unsupported database system",
				"DatabaseConfig port is required",
				"DatabaseConfig root password is required",
			},
			actualErrs:   errs,
			validatedVal: db,
		}

		failTestOnUnspottedError(res, t)
	} else {
		t.Errorf("Did not return any errors for value %v", db)
	}
}

func TestDatabase_ValidateCorrectInput(t *testing.T) {
	db := service.DatabaseConfig{
		System: service.MySQL,
		Port:   3306,
		Name:   "testdb",
		Credentials: service.Credentials{
			RootPassword: "rootpass",
		},
	}

	errs := db.Validate()

	failTestOnErrorsOnCorrectInput(errs, t)
}
