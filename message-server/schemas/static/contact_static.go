package static

import (
	"context"
	"sync"

	"github.com/kaleb-white/letthemknow/message-server/db"
	"github.com/kaleb-white/letthemknow/message-server/schemas/contact"
)

const SOURCE_CONTACT_STATIC = "contact_static"

var once sync.Once

func InitializeContacts(ctx context.Context) error {
	var err error
	once.Do(func () {
		err = db.InitializeTable(ctx, contact.CONTACT_TABLEDEF, "contact table")
	})
	if err != nil {
		return err
	}

	return nil
}
