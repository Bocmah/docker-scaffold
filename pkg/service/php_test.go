package service_test

import (
	"reflect"
	"testing"

	"github.com/Bocmah/phpdocker-gen/pkg/service"
)

func TestPHP_FillDefaultsIfNotSet(t *testing.T) {
	php := service.PHPConfig{}

	php.FillDefaultsIfNotSet()

	want := service.PHPConfig{
		Version:    "7.4",
		Extensions: []string{"mbstring", "zip", "exif", "pcntl", "gd"},
	}

	if !reflect.DeepEqual(php, want) {
		t.Errorf("PHPConfig FillDefaultsIfNotSet incorrect defaults. Want %v, got %v", want, php)
	}
}

func TestPHP_AddDatabaseExtension(t *testing.T) {
	php := service.PHPConfig{Version: "7.4", Extensions: []string{}}

	php.AddDatabaseExtension(service.MySQL)

	wantExt := "pdo_mysql"
	hasExt := false
	for _, ext := range php.Extensions {
		if ext == wantExt {
			hasExt = true
			break
		}
	}

	if !hasExt {
		t.Errorf("Failed to add extension for MySQL")
	}

	php.AddDatabaseExtension(service.PostgreSQL)

	wantExt = "pdo_pgsql"
	hasExt = false
	for _, ext := range php.Extensions {
		if ext == wantExt {
			hasExt = true
			break
		}
	}

	if !hasExt {
		t.Errorf("Failed to add extension for PostgreSQL")
	}
}

func TestPHP_ValidateIncorrectInput(t *testing.T) {
	php := service.PHPConfig{}

	errs := php.Validate()

	if errs != nil {
		res := validationResult{
			wantErrs:     []string{"PHPConfig version is required"},
			actualErrs:   errs,
			validatedVal: php,
		}

		failTestOnUnspottedError(res, t)
	} else {
		t.Errorf("Did not return any errors for value %v", php)
	}
}

func TestPHP_ValidateCorrectInput(t *testing.T) {
	php := service.PHPConfig{Version: "7.4"}

	errs := php.Validate()

	failTestOnErrorsOnCorrectInput(errs, t)
}

func TestPHP_IsEmpty(t *testing.T) {
	if !(&service.PHPConfig{}).IsEmpty() {
		t.Errorf("Failed to assert that empty PHPConfig service is actually empty")
	}

	if (&service.PHPConfig{Version: "7.4"}).IsEmpty() {
		t.Errorf("Failed to assert that non-empty PHPConfig service is actually non-empty")
	}
}
