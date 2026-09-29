package static

import (
	"context"
	"sync"

	"github.com/kaleb-white/letthemknow/message-server/db"
	cl "github.com/kaleb-white/letthemknow/message-server/schemas/contact_list"
)

const SOURCE_CONTACT_LIST_STATIC = "contact_list_static"

var onceContactList sync.Once

func InitializeContactLists(ctx context.Context) error {
	var err error
	onceContactList.Do(func () {
		err = db.InitializeTable(ctx, cl.CONTACT_LIST_TABLEDEF, "contact list table")
	})
	if err != nil {
		return err
	}

	return nil
}

