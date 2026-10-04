package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/kaleb-white/letthemknow/message-server/models"
	"github.com/kaleb-white/letthemknow/message-server/sqlite"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)


func TestSqliteContactList(t *testing.T) {
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

	contactListStore := sqlite.ContactListStore{}

	var id uint64
	var ctx = context.Background()

	t.Run("TestInitializeContactDetailWithContact", func(t *testing.T) {
		err := contactListStore.Init(ctx)
		if err != nil {
			t.Errorf("Error during initialization: %s", err.Error())
		}
		localId, err := contactListStore.Write(ctx, &contactList)
		id = localId
		if err != nil {
			t.Errorf("Error during write of base contact list: %s", err.Error())
		}
	})

	t.Run("TestReadContactDetailGoldenPath", func(t *testing.T) {
		copyContactList, err := contactListStore.Read(ctx, id)
		if err != nil {
			t.Errorf("Error while retrieving contact list store: %s", err.Error())
		}
		errs, equals := copyContactList.Equals(&contactList)
		if !equals {
			t.Errorf("ContactLists were not equal: %s", utils.PrettyPrintErrors(errs))
		}
	})

	t.Run("TestDeleteContactDetailGoldenPath", func(t *testing.T) {
		localId, err := contactListStore.Delete(ctx, contactList.Id)
		if err != nil {
			t.Errorf("Failed to delete cd: %s", err.Error())
		}
		if localId != id {
			t.Errorf("Deleted id did not equal id, was: %d", localId)
		}
	})
}
