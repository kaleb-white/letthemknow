package cd_test

import (
	"context"
	"testing"
	"time"

	"github.com/kaleb-white/letthemknow/message-server/schemas/contact_detail"
	"github.com/kaleb-white/letthemknow/message-server/schemas/contact"
	"github.com/kaleb-white/letthemknow/message-server/schemas/static"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)



var baseContactDetail cd.ContactDetail = cd.ContactDetail{
	Id: 1,
	Phone: 2223518234,
	FirstName: "John",
	LastName: "Doe",
	LastUpdatedAtDateTime: time.Now(),
	LastUpdatedBy: "test",
	Org: "test",
}

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

func TestValidateContactDetailValidSchema(t *testing.T) {
	err, errs := baseContactDetail.Validate()
	if err {
		t.Errorf("Valid Schema should pass validation. %s", utils.PrettyPrintErrors(errs))
	}
}

func TestValidateContactDetailValidSchemaWithoutField(t *testing.T) {
	baseContactDetail.Org = ""
	err, errs := baseContactDetail.Validate()
	if err {
		t.Errorf("Valid Schema minus one non-allowed field should pass validation. %s", utils.PrettyPrintErrors(errs))
	}
	baseContactDetail.Org = "test"
}

func TestValidateContactDetailMissingField(t *testing.T) {
	baseContactDetail.Phone = 0
	err, _ := baseContactDetail.Validate()
	if !err {
		t.Errorf("Invalid schema (missing Phone) should not pass validation.")
	}
	baseContactDetail.Phone = 2223518234
}

func TestValidateContactDetailMissingSeveralFields(t *testing.T) {
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
}

var id uint64
var ctx = context.Background()

func TestInitializeContactDetailWithContact(t* testing.T) {
	err := static.InitializeContacts(ctx)
	if err != nil {
		t.Errorf("Error during initialization: %s", err.Error())
	}
	localId, err := baseContact.Write(ctx)
	id = localId
	if err != nil {
		t.Errorf("Error during write of base contact: %s", err.Error())
	}
}

func TestReadContactDetailGoldenPath(t* testing.T) { 
	copyContactDetail := cd.ContactDetail{}		
	copyContactDetail.Read(ctx, id)
	errs, equals := copyContactDetail.Equals(&baseContactDetail)
	if !equals {
		t.Errorf("ContactDetails were not equal: %s", utils.PrettyPrintErrors(errs))
	}
}

func TestDeleteContactDetailGoldenPath(t* testing.T) {
	localId, err := baseContactDetail.Delete(ctx)
	if err != nil {
		t.Errorf("Failed to delete cd: %s", err.Error())
	}
	if localId != id {
		t.Errorf("Deleted id did not equal id, was: %d", localId)
	}
}
