package service

import "strings"

// ValidationErrors is a collection of validation errors
type ValidationErrors []string

func (v ValidationErrors) Error() string {
	return strings.Join(v, "\n")
}

// Add adds error/errors to the collection
func (v *ValidationErrors) Add(err ...string) {
	*v = append(*v, err...)
}

// IsEmpty determines whether collection is empty
func (v ValidationErrors) IsEmpty() bool {
	return len(v) == 0
}

// Has determines whether collection has given error
func (v ValidationErrors) Has(err string) bool {
	for _, el := range v {
		if el == err {
			return true
		}
	}

	return false
}

// Merge merges current collection of errors with another collection of errors
func (v *ValidationErrors) Merge(errs *ValidationErrors) {
	v.Add(*errs...)
}
