package schemas

type Savable interface {
	Read(Id string) (*Savable, error)
	Write() (string, error)
	Delete(Id string) error
}

type Validatable interface {
	// First rv true if there was an error, second is all the errors
	Validate() (bool, []error)
}
