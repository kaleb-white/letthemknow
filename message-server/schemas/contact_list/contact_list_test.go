package cl_test

import (
	"time"
	"testing"
	"context"

	cl "github.com/kaleb-white/letthemknow/message-server/schemas/contact_list"
	"github.com/kaleb-white/letthemknow/message-server/utils"
	"github.com/kaleb-white/letthemknow/message-server/schemas/static"

)

var contactList cl.ContactList = cl.ContactList{
	Id: 1,	
	Name: "MyList",
	Description: "123",
	ContactList: []uint64{1, 2, 3},
	CreatedAtDateTime: time.Now(),
	CreatedBy: "test",
	LastUpdatedAtDateTime: time.Now(),
	LastUpdatedBy: "test2",
}

func TestValidateContactListValidSchema(t *testing.T) {
	err, errs := contactList.Validate()
	if err {
		t.Errorf("Valid Schema should pass validation. %s", utils.PrettyPrintErrors(errs))
	}
}

func TestValidateContactDetailValidSchemaWithoutField(t *testing.T) {
	contactList.Description = ""
	defer func() {
		contactList.Description = "123"
	}()
	err, errs := contactList.Validate()
	if err {
		t.Errorf("Valid Schema minus one non-allowed field should pass validation. %s", utils.PrettyPrintErrors(errs))
	}
	contactList.Description = "123"
}

func TestValidateContactDetailMissingField(t *testing.T) {
	contactList.Name = ""
	defer func() {
		contactList.Name = "MyList"
	}()
	err, _ := contactList.Validate()
	if !err {
		t.Errorf("Invalid schema (missing Phone) should not pass validation.")
	}
}

func TestValidateContactDetailMissingSeveralFields(t *testing.T) {
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
}

var id uint64
var ctx = context.Background()

func TestInitializeContactDetailWithContact(t *testing.T) {
	err := static.InitializeContactLists(ctx)
	if err != nil {
		t.Errorf("Error during initialization: %s", err.Error())
	}
	localId, err := contactList.Write(ctx)
	id = localId
	if err != nil {
		t.Errorf("Error during write of base contact list: %s", err.Error())
	}
}

func TestReadContactDetailGoldenPath(t *testing.T) {
	copyContactList := cl.ContactList{}
	copyContactList.Read(ctx, id)
	errs, equals := copyContactList.Equals(&contactList)
	if !equals {
		t.Errorf("ContactDetails were not equal: %s", utils.PrettyPrintErrors(errs))
	}
}

func TestDeleteContactDetailGoldenPath(t *testing.T) {
	localId, err := contactList.Delete(ctx)
	if err != nil {
		t.Errorf("Failed to delete cd: %s", err.Error())
	}
	if localId != id {
		t.Errorf("Deleted id did not equal id, was: %d", localId)
	}
}
