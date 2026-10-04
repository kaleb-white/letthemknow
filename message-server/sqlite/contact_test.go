package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/kaleb-white/letthemknow/message-server/models"
	"github.com/kaleb-white/letthemknow/message-server/sqlite"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)

func TestContactStore(t *testing.T) {
	baseContact := models.Contact{
		Id: 1,
		Phone: 2223518234,
		Phone2: 2223518234,
		Phone3: 2223518234,
		FirstName: "John",
		LastName: "Doe",
		FullName: "John D. Doe",
		ListMembership: []uint64{1, 2,},
		CreatedAtDateTime: time.Now(),
		CreatedBy: "test",
		LastUpdatedAtDateTime: time.Now(),
		LastUpdatedBy: "test",
		Org: "test",
	}

	contactStore := sqlite.ContactStore{}

	var id uint64
	var ctx = context.Background()

	t.Run("TestInitializeContactTable", func(t* testing.T) {
		err := contactStore.Init(ctx)
		if err != nil {
			t.Errorf("Error during initialization: %s", err.Error())
		}
	})

	t.Run("TestWriteContactGoldenPath", func(t* testing.T) {
		localId, err := contactStore.Write(ctx, &baseContact)
		id = localId
		if err != nil {
			t.Errorf("Error during write: %s", err.Error())
		}
	})

	t.Run("TestReadContactGoldenPath", func(t* testing.T) { 
		copyContact, err := contactStore.Read(ctx, id)
		if err != nil {
			t.Errorf("Failed to get contact to test equals: %s", err.Error())
		}
		errs, equals := copyContact.Equals(&baseContact)
		if !equals {
			t.Errorf("Contacts were not equal: %s", utils.PrettyPrintErrors(errs))
		}
	})

	t.Run("TestReadContactListMembershipRebuilt", func(t* testing.T) { 
		copyContact, err := contactStore.Read(ctx, id)
		if err != nil {
			t.Errorf("Failed to get contact to test equals: %s", err.Error())
		}
		membershipLen := len(copyContact.ListMembership)
		if membershipLen != 2 {
			t.Errorf("Expected membership list to have 2 ids, had %d", membershipLen)
		}
		expect1 := copyContact.ListMembership[0]
		expect2 := copyContact.ListMembership[1]
		if expect1 != 1 || expect2 != 2 {
			t.Errorf("Unmarshaling failed, expected 1 and 2 and got: %d, %d")
		} 
	})

	t.Run("TestDeleteContactGoldenPath", func(t* testing.T) {
		localId, err := contactStore.Delete(ctx, baseContact.Id)
		if err != nil {
			t.Errorf("Failed to delete contact: %s", err.Error())
		}
		if localId != id {
			t.Errorf("Deleted id did not equal id, was: %d", localId)
		}
	})
}
