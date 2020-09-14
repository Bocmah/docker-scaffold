package service_test

import (
	"testing"

	"github.com/Bocmah/phpdocker-gen/pkg/service"
)

func equal(e1 service.ValidationErrors, e2 service.ValidationErrors) bool {
	if len([]string(e1)) != len([]string(e2)) {
		return false
	}

	for i, v := range e1 {
		if v != e2[i] {
			return false
		}
	}
	return true
}

func TestValidationErrorsErrorString(t *testing.T) {
	errors := service.ValidationErrors{"Sample error", "Sample error 2"}

	got := errors.Error()
	want := "Sample error\nSample error 2"

	if got != want {
		t.Errorf("Incorrect error string. Want %v. Got %v", want, got)
	}
}

func TestValidationErrors_AddAdd(t *testing.T) {
	got := service.ValidationErrors{}

	got.Add("Sample error")
	got.Add("Sample error 2")

	want := service.ValidationErrors{"Sample error", "Sample error 2"}

	if !equal(got, want) {
		t.Errorf("Failed to add error. Want %v. Got %v", want, got)
	}
}

func TestValidationErrorsEmptyAdd(t *testing.T) {
	got := service.ValidationErrors{}

	got.Add("Sample error")
	got.Add()

	want := service.ValidationErrors{"Sample error"}

	if !equal(got, want) {
		t.Errorf("Empty add affects contents. Want %v. Got %v", want, got)
	}
}

func TestValidationErrors_IsEmpty(t *testing.T) {
	if !(&service.ValidationErrors{}).IsEmpty() {
		t.Errorf("Failed to validate that empty ValidationErrors is actually empty")
	}

	if (&service.ValidationErrors{"Not empty"}).IsEmpty() {
		t.Errorf("Failed to validate that non empty ValidationErrors is actually non empty")
	}
}

func TestValidationErrors_Merge(t *testing.T) {
	got := &service.ValidationErrors{"Error 1"}

	got.Merge(&service.ValidationErrors{"Error 2"})
	want := service.ValidationErrors{"Error 1", "Error 2"}

	if !equal(*got, want) {
		t.Errorf("Incorrect merge. Want %v. Got %v", want, got)
	}
}

func TestValidationErrors_Has(t *testing.T) {
	tests := map[string]struct {
		input service.ValidationErrors
		error string
		want  bool
	}{
		"simple": {
			input: service.ValidationErrors{"Test error"},
			error: "Test error",
			want:  true,
		},
		"more than one error": {
			input: service.ValidationErrors{"Test error", "Another error"},
			error: "Test error",
			want:  true,
		},
		"doesn't have": {
			input: service.ValidationErrors{"Test error"},
			error: "Another error",
			want:  false,
		},
		"empty": {
			input: service.ValidationErrors{},
			error: "Test error",
			want:  false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.Has(tc.error)

			if got != tc.want {
				t.Fatalf("incorrect ValidationErrors.Has() behaviour. got %v want %v", got, tc.want)
			}
		})
	}
}
