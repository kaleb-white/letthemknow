package models_test

import (
	"testing"
	"time"

	"github.com/kaleb-white/letthemknow/message-server/models"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)

func TestContact(t *testing.T) {
	baseContact := models.Contact{
		Id:                    1,
		Phone:                 2223518234,
		Phone2:                2223518234,
		Phone3:                2223518234,
		FirstName:             "John",
		LastName:              "Doe",
		FullName:              "John D. Doe",
		CreatedAtDateTime:     time.Now(),
		CreatedBy:             "test",
		LastUpdatedAtDateTime: time.Now(),
		LastUpdatedBy:         "test",
		Org:                   "test",
	}
	t.Run("TestValidateContactValidSchema", func(t *testing.T) {
		err, errs := baseContact.Validate()
		if err {
			t.Errorf("Valid Schema should pass validation. %s", utils.PrettyPrintErrors(errs))
		}
	})

	t.Run("TestValidateContactValidSchemaWithoutField", func(t *testing.T) {
		baseContact.Org = ""
		err, errs := baseContact.Validate()
		if err {
			t.Errorf("Valid Schema minus one non-allowed field should pass validation. %s", utils.PrettyPrintErrors(errs))
		}
		baseContact.Org = "test"
	})

	t.Run("TestValidateContactMissingField", func(t *testing.T) {
		baseContact.Phone = 0
		err, _ := baseContact.Validate()
		if !err {
			t.Errorf("Invalid schema (missing Phone) should not pass validation.")
		}
		baseContact.Phone = 2223518234
	})

	t.Run("TestValidateContactMissingSeveralFields", func(t *testing.T) {
		baseContact.Phone = 0
		baseContact.FirstName = ""
		err, errs := baseContact.Validate()
		if !err {
			t.Errorf("Invalid schema (missing Phone, Id) should not pass validation.")
		}
		if len(errs) != 2 {
			t.Errorf("Expected two errors. Got %d. %s", len(errs), utils.PrettyPrintErrors(errs))
		}
		baseContact.FirstName = "John"
		baseContact.Phone = 2223518234
	})
}
