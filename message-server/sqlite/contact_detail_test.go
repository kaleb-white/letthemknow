package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/kaleb-white/letthemknow/message-server/models"
	"github.com/kaleb-white/letthemknow/message-server/sqlite"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)

func TestContactDetailStore(t *testing.T) {
	contactDetailStore := sqlite.ContactDetailStore{}

	baseContactDetail := models.ContactDetail{
		Id: 1,
		Phone: 2223518234,
		FirstName: "John",
		LastName: "Doe",
		LastUpdatedAtDateTime: time.Now(),
		LastUpdatedBy: "test",
		Org: "test",
	}

	contactStore := sqlite.ContactStore{}

	baseContact := models.Contact{
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


	var id uint64
	var ctx = context.Background()

	t.Run("TestInitializeContactDetailWithContact", func(t* testing.T) {
		err := contactDetailStore.Init(ctx)
		if err != nil {
			t.Errorf("Error during initialization: %s", err.Error())
		}
		localId, err := contactStore.Write(ctx, &baseContact)
		id = localId
		if err != nil {
			t.Errorf("Error during write of base contact: %s", err.Error())
		}
	})

	t.Run("TestReadContactDetailGoldenPath", func(t* testing.T) { 
		copyContactDetail, _ := contactDetailStore.Read(ctx, baseContactDetail.Id)	
		errs, equals := copyContactDetail.Equals(&baseContactDetail)
		if !equals {
			t.Errorf("ContactDetails were not equal: %s", utils.PrettyPrintErrors(errs))
		}
	})

	t.Run("TestDeleteContactDetailGoldenPath", func(t* testing.T) {
		localId, err := contactDetailStore.Delete(ctx, baseContactDetail.Id)
		if err != nil {
			t.Errorf("Failed to delete cd: %s", err.Error())
		}
		if localId != id {
			t.Errorf("Deleted id did not equal id, was: %d", localId)
		}
	})
}
