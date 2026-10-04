package models_test

import (
	"testing"
	"time"

	"github.com/kaleb-white/letthemknow/message-server/models"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)

func TestContactList(t *testing.T) {
	contactList := models.ContactList{
		Id: 1,	
		Name: "MyList",
		Description: "123",
		ContactList: []uint64{1, 2, 3},
		CreatedAtDateTime: time.Now(),
		CreatedBy: "test",
		LastUpdatedAtDateTime: time.Now(),
		LastUpdatedBy: "test2",
	}

	t.Run("TestValidateContactListValidSchema", func(t *testing.T) {
		err, errs := contactList.Validate()
		if err {
			t.Errorf("Valid Schema should pass validation. %s", utils.PrettyPrintErrors(errs))
		}
	})

	t.Run("TestValidateContactDetailValidSchemaWithoutField", func(t *testing.T) {
		contactList.Description = ""
		defer func() {
			contactList.Description = "123"
		}()
		err, errs := contactList.Validate()
		if err {
			t.Errorf("Valid Schema minus one non-allowed field should pass validation. %s", utils.PrettyPrintErrors(errs))
		}
		contactList.Description = "123"
	})

	t.Run("TestValidateContactDetailMissingField", func(t *testing.T) {
		contactList.Name = ""
		defer func() {
			contactList.Name = "MyList"
		}()
		err, _ := contactList.Validate()
		if !err {
			t.Errorf("Invalid schema (missing Phone) should not pass validation.")
		}
	})

	t.Run("TestValidateContactDetailMissingSeveralFields", func(t *testing.T) {
		contactList.Name = ""
		contactList.CreatedBy = ""
		defer func() {
			contactList.Name = "MyList"
			contactList.CreatedBy = "test"
		}()
		err, errs := contactList.Validate()
		if !err {
			t.Errorf("Invalid schema (missing Name, CreatedBy) should not pass validation.")
		}
		if len(errs) != 2 {
			t.Errorf("Expected two errors. Got %d. %s", len(errs), utils.PrettyPrintErrors(errs))
		}
	})
}
