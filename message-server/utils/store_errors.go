package utils

type StoreError interface {
	Error() string
	Status() int
}

type storeError struct {
	err string
	status int
}

func (e storeError) Error() string {
	return e.err
}

func (e storeError) Status() int {
	return e.status
}

func NewStoreError(err string, status int) StoreError {
	return storeError{
		err: err,
		status: status,
	}	
}
