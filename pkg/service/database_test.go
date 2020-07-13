package service_test

import (
	"github.com/Bocmah/phpdocker-scaffold/pkg/service"
	"testing"
)

type validationResult struct {
	wantErrs     []string
	actualErrs   service.ValidationErrors
	validatedVal interface{}
}

func failTestOnUnspottedError(result validationResult, t *testing.T) {
	for _, e := range result.wantErrs {
		if !result.actualErrs.Has(e) {
			t.Errorf("Failed to spot error %s in value %v", e, result.validatedVal)
		}
	}
}

func failTestOnErrorsOnCorrectInput(errs *service.ValidationErrors, t *testing.T) {
	if errs != nil {
		t.Errorf("Following errors were returned despite correct inputs %v", *errs)
	}
}

func TestDatabase_FillDefaultsIfNotSet(t *testing.T) {
	db := service.Database{}

	db.FillDefaultsIfNotSet()

	want := service.Database{
		System:  service.MySQL,
		Port:    3306,
		Version: "8.0",
	}

	if db != want {
		t.Errorf("Incorrect defaults, want %v, got %v", want, db)
	}
}

func TestDatabase_ValidateIncorrectInput(t *testing.T) {
	db := service.Database{System: "Unsupported"}

	errs := db.Validate()

	if errs != nil {
		res := validationResult{
			wantErrs: []string{
				"Unsupported database system",
				"Database name is required",
				"Database port is required",
				"Database root password is required",
			},
			actualErrs: *errs,
			validatedVal: db,
		}

		failTestOnUnspottedError(res, t)
	} else {
		t.Errorf("Did not return any errors for value %v", db)
	}
}

func TestDatabase_ValidateCorrectInput(t *testing.T) {
	db := service.Database{
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
