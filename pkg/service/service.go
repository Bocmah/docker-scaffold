package service

type Service interface {
	FillDefaultsIfNotSet()
	Validate() error
}
