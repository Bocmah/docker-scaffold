package service_test

import (
	"testing"

	"github.com/Bocmah/phpdocker-scaffold/pkg/service"
)

func TestNodeJS_FillDefaultsIfNotSet(t *testing.T) {
	nodejs := service.NodeJSConfig{}

	nodejs.FillDefaultsIfNotSet()

	want := service.NodeJSConfig{
		Version: "latest",
	}

	if nodejs != want {
		t.Errorf("Incorrect defaults, want %v, got %v", want, nodejs)
	}
}

func TestNodeJS_ValidateIncorrectInput(t *testing.T) {
	nodejs := service.NodeJSConfig{}

	errs := nodejs.Validate()

	if errs != nil {
		res := validationResult{
			wantErrs: []string{
				"Node.js version is required",
			},
			actualErrs:   errs,
			validatedVal: nodejs,
		}

		failTestOnUnspottedError(res, t)
	} else {
		t.Errorf("Did not return any errors for value %v", nodejs)
	}
}

func TestNodeJS_ValidateCorrectInput(t *testing.T) {
	nodejs := service.NodeJSConfig{Version: "latest"}

	errs := nodejs.Validate()

	failTestOnErrorsOnCorrectInput(errs, t)
}
