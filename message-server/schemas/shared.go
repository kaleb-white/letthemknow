package schemas

import "context"

type Savable interface {
	Initialize(ctx context.Context) error
	Read(ctx context.Context, Id uint64) (*Savable, error)
	Write(ctx context.Context) (uint64, error)
	Delete(ctx context.Context) (string, error)
}

type Validatable interface {
	// First rv true if there was an error, second is all the errors
	Validate() (bool, []error)
}
