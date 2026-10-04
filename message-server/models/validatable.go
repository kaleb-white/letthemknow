package models

// First rv true if there was an error, second is all the errors
type Validatable interface {
	Validate() (bool, []error)
}
