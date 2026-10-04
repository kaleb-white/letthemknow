package models

type IdentifiableValidatable[M any] interface {
	*M
	Validatable
	Identifiable
}
