package handlers

import (
	"net/http"
	"time"

	"github.com/kaleb-white/letthemknow/message-server/models"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)

const SOURCE_CONTACT_HANDLERS string = "contact_handler"

func HandleReadContact(contactStore *models.ContactStore, getenv func(string) string) http.Handler {
	return genericReadHandler(contactStore, getenv, SOURCE_CONTACT_HANDLERS, "Contact")
}

func HandleWriteContact(contactStore *models.ContactStore, getenv func(string) string) http.Handler {
	return genericWriteHandler(
		contactStore, 
		getenv, 
		SOURCE_CONTACT_HANDLERS, 
		"Contact",
		func(c *models.Contact, r *http.Request) {
			// Set before validation
			c.CreatedBy = utils.GetAttribution(r.Context())
			c.CreatedAtDateTime = time.Now()			
			c.LastUpdatedBy = utils.GetAttribution(r.Context())
			c.LastUpdatedAtDateTime = time.Now()
		},
	)
}

func HandleWriteContactById(contactStore *models.ContactStore, getenv func(string) string) http.Handler {
	return genericWriteByIdHandler(
		contactStore, 
		getenv, 
		SOURCE_CONTACT_HANDLERS, 
		"Contact",
		func(c *models.Contact, r *http.Request) {
			// Set before validation
			c.CreatedBy = utils.GetAttribution(r.Context())
			c.CreatedAtDateTime = time.Now()			
			c.LastUpdatedBy = utils.GetAttribution(r.Context())
			c.LastUpdatedAtDateTime = time.Now()
		},
	)

}

func HandleDeleteContact(contactStore *models.ContactStore, getenv func(string) string) http.Handler {
	return genericDeleteHandler[models.Contact](contactStore, getenv, SOURCE_CONTACT_HANDLERS, "Contact")
}
