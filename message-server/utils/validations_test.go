package utils_test

import (
	"testing"

	"github.com/kaleb-white/letthemknow/message-server/schemas/contact"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)

func CheckRequiredFieldsArentDefaultNoErrors(t *testing.T) {
	c := contact.Contact{
		Id: 123,
		Phone: 123,
	}

	errs := []error{}
	fields := []string{"Id", "Phone"}

	utils.CheckRequiredFieldsArentDefault(&c, &fields, &errs)

	if len(errs) > 0 {
		t.Errorf("Should've gotten 0 errors, received more than 0")
	}	
}

func CheckRequiredFieldsArentDefaultSomeErrors(t *testing.T) {
	c := contact.Contact{
		Id: 123,
		Phone: 123,
	}

	errs := []error{}
	fields := []string{"Id", "Phone", "Phone1"}

	utils.CheckRequiredFieldsArentDefault(&c, &fields, &errs)

	if len(errs) != 1 {
		t.Errorf("Should've gotten 1 error, received %d.", len(errs))
	}	
}

func CheckRequiredFieldsArentDefaultExpectPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Should've panic, defer fired without error to recover")
		}
	}()

	c := contact.Contact{
		Id: 123,
		Phone: 123,
	}

	errs := []error{}
	fields := []string{"Id", "Phone", "Phone1", "aaabb"}

	utils.CheckRequiredFieldsArentDefault(&c, &fields, &errs)

	if len(errs) != 1 {
		t.Errorf("Should've gotten 1 error, received %d.", len(errs))
	}	
}
