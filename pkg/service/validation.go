package service

import "strings"

type ValidationErrors []string

func (v ValidationErrors) Error() string {
	return strings.Join(v, "\n")
}

func (v *ValidationErrors) Add(err ...string) {
	*v = append(*v, err...)
}

func (v ValidationErrors) IsEmpty() bool {
	return len(v) == 0
}

func (v ValidationErrors) Has(err string) bool {
	for _, el := range v {
		if el == err {
			return true
		}
	}

	return false
}

func (v *ValidationErrors) Merge(errs *ValidationErrors) {
	v.Add(*errs...)
}
