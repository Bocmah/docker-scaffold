package service_test

import (
	"github.com/Bocmah/phpdocker-scaffold/pkg/service"
	"testing"
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

func TestErrorString(t *testing.T) {
	errors := service.ValidationErrors{"Sample error", "Sample error 2"}

	got := errors.Error()
	want := "Sample error\nSample error 2"

	if got != want {
		t.Errorf("Incorrect error string. Want %v. Got %v", want, got)
	}
}

func TestAdd(t *testing.T) {
	got := service.ValidationErrors{}

	got.Add("Sample error")
	got.Add("Sample error 2")

	want := service.ValidationErrors{"Sample error", "Sample error 2"}

	if !equal(got, want) {
		t.Errorf("Failed to add error. Want %v. Got %v", want, got)
	}
}

func TestIsEmpty(t *testing.T) {
	if !(&service.ValidationErrors{}).IsEmpty() {
		t.Errorf("Failed to validate that empty ValidationErrors is actually empty")
	}

	if (&service.ValidationErrors{"Not empty"}).IsEmpty() {
		t.Errorf("Failed to validate that non empty ValidationErrors is actually non empty")
	}
}

func TestMerge(t *testing.T) {
	got := &service.ValidationErrors{"Error 1"}

	got.Merge(&service.ValidationErrors{"Error 2"})
	want := service.ValidationErrors{"Error 1", "Error 2"}

	if !equal(*got, want) {
		t.Errorf("Incorrect merge. Want %v. Got %v", want, got)
	}
}