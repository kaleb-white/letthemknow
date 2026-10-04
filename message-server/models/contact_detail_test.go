package models_test

import (
	"testing"
	"time"

	"github.com/kaleb-white/letthemknow/message-server/models"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)

func TestContactDetail(t *testing.T) {
	baseContactDetail := models.ContactDetail{
		Id:                    1,
		Phone:                 2223518234,
		FirstName:             "John",
		LastName:              "Doe",
		LastUpdatedAtDateTime: time.Now(),
		LastUpdatedBy:         "test",
		Org:                   "test",
	}

	t.Run("TestValidateContactDetailValidSchema", func(t *testing.T) {
		err, errs := baseContactDetail.Validate()
		if err {
			t.Errorf("Valid Schema should pass validation. %s", utils.PrettyPrintErrors(errs))
		}
	})

	t.Run("TestValidateContactDetailValidSchemaWithoutField", func(t *testing.T) {
		baseContactDetail.Org = ""
		err, errs := baseContactDetail.Validate()
		if err {
			t.Errorf("Valid Schema minus one non-allowed field should pass validation. %s", utils.PrettyPrintErrors(errs))
		}
		baseContactDetail.Org = "test"
	})

	t.Run("TestValidateContactDetailMissingField", func(t *testing.T) {
		baseContactDetail.Phone = 0
		err, _ := baseContactDetail.Validate()
		if !err {
			t.Errorf("Invalid schema (missing Phone) should not pass validation.")
		}
		baseContactDetail.Phone = 2223518234
	})

	t.Run("TestValidateContactDetailMissingSeveralFields", func(t *testing.T) {
		baseContactDetail.Phone = 0
		baseContactDetail.FirstName = ""
		err, errs := baseContactDetail.Validate()
		if !err {
			t.Errorf("Invalid schema (missing Phone, Id) should not pass validation.")
		}
		if len(errs) != 2 {
			t.Errorf("Expected two errors. Got %d. %s", len(errs), utils.PrettyPrintErrors(errs))
		}
		baseContactDetail.FirstName = "John"
		baseContactDetail.Phone = 2223518234
	})
}
