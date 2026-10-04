package handlers

import (
	"net/http"
	"time"

	"github.com/kaleb-white/letthemknow/message-server/models"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)

const SOURCE_CONTACT_LIST_HANDLER string = "contact_list_handler"

func HandleReadContactList(contactListStore *models.ContactListStore, getenv func(string) string) http.Handler {
	return genericReadHandler(contactListStore, getenv, SOURCE_CONTACT_LIST_HANDLER, "ContactList")
}

func HandleWriteContactList(contactListStore *models.ContactListStore, getenv func(string) string) http.Handler {
	return genericWriteHandler(
		contactListStore, 
		getenv, 
		SOURCE_CONTACT_LIST_HANDLER, 
		"ContactList", 
		func(c *models.ContactList, r *http.Request) {
			c.CreatedBy = utils.GetAttribution(r.Context())
			c.CreatedAtDateTime = time.Now()			
			c.LastUpdatedBy = utils.GetAttribution(r.Context())
			c.LastUpdatedAtDateTime = time.Now()
	})
}

func HandleWriteContactListById(contactListStore *models.ContactListStore, getenv func(string) string) http.Handler {
	return genericWriteByIdHandler(
		contactListStore, 
		getenv, 
		SOURCE_CONTACT_LIST_HANDLER, 
		"ContactList", 
		func(c *models.ContactList, r *http.Request) {
			c.CreatedBy = utils.GetAttribution(r.Context())
			c.CreatedAtDateTime = time.Now()			
			c.LastUpdatedBy = utils.GetAttribution(r.Context())
			c.LastUpdatedAtDateTime = time.Now()
	})
}

func HandleDeleteContactList(contactListStore *models.ContactListStore, getenv func(string) string) http.Handler {
	return genericDeleteHandler[models.Contact](contactListStore, getenv, SOURCE_CONTACT_HANDLERS, "Contact")
}

