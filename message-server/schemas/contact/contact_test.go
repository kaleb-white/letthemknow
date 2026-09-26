package contact_test

import (
	"context"
	"testing"
	"time"

	"github.com/kaleb-white/letthemknow/message-server/schemas/contact"
	"github.com/kaleb-white/letthemknow/message-server/schemas/static"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)

var baseContact contact.Contact = contact.Contact{
	Id: 1,
	Phone: 2223518234,
	Phone2: 2223518234,
	Phone3: 2223518234,
	FirstName: "John",
	LastName: "Doe",
	FullName: "John D. Doe",
	CreatedAtDateTime: time.Now(),
	CreatedBy: "test",
	LastUpdatedAtDateTime: time.Now(),
	LastUpdatedBy: "test",
	Org: "test",
}

func TestValidateContactValidSchema(t *testing.T) {
	err, errs := baseContact.Validate()
	if err {
		t.Errorf("Valid Schema should pass validation. %s", utils.PrettyPrintErrors(errs))
	}
}

func TestValidateContactValidSchemaWithoutField(t *testing.T) {
	baseContact.Org = ""
	err, errs := baseContact.Validate()
	if err {
		t.Errorf("Valid Schema minus one non-allowed field should pass validation. %s", utils.PrettyPrintErrors(errs))
	}
	baseContact.Org = "test"
}

func TestValidateContactMissingField(t *testing.T) {
	baseContact.Phone = 0
	err, _ := baseContact.Validate()
	if !err {
		t.Errorf("Invalid schema (missing Phone) should not pass validation.")
	}
	baseContact.Phone = 2223518234
}

func TestValidateContactMissingSeveralFields(t *testing.T) {
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
}

var id uint64
var ctx = context.Background()

func TestInitializeContactTable(t* testing.T) {
	err := static.InitializeContacts(ctx)
	if err != nil {
		t.Errorf("Error during initialization: %s", err.Error())
	}
}

func TestWriteContactGoldenPath(t* testing.T) {
	localId, err := baseContact.Write(ctx)
	id = localId
	if err != nil {
		t.Errorf("Error during write: %s", err.Error())
	}
}

func TestReadContactGoldenPath(t* testing.T) { 
	copyContact := contact.Contact{}		
	copyContact.Read(ctx, id)
	errs, equals := copyContact.Equals(&baseContact)
	if !equals {
		t.Errorf("Contacts were not equal: %s", utils.PrettyPrintErrors(errs))
	}
}

func TestDeleteContactGoldenPath(t* testing.T) {
	localId, err := baseContact.Delete(ctx)
	if err != nil {
		t.Errorf("Failed to delete contact: %s", err.Error())
	}
	if localId != id {
		t.Errorf("Deleted id did not equal id, was: %d", localId)
	}
}
