package main

import (
	"net/http"
	"os"

	"github.com/kaleb-white/letthemknow/message-server/handlers"
	"github.com/kaleb-white/letthemknow/message-server/middleware"
	"github.com/kaleb-white/letthemknow/message-server/models"
)

var getenv func(string) string = os.Getenv

func AddRoutes(
	mux *http.ServeMux,
	config *models.Config,
) {
	
	/* ----- Contacts ----- */
	// Read
	mux.Handle("GET /api/v1/contacts/{id}",
		middleware.HasSlugs(
			handlers.HandleReadContact(
				config.ContactStore, 
				getenv,
			),
			[]string{"id"},
		),
	)

	// Write
	mux.Handle("POST /api/v1/contacts",
		handlers.HandleWriteContact(
			config.ContactStore,
			getenv,
		),
	)

	// Update
	mux.Handle("POST /api/v1/contacts/{id}",
		middleware.HasSlugs(
			handlers.HandleWriteContactById(
				config.ContactStore, 
				getenv,
			),
			[]string{"id"},
		),
	)

 // Delete	
 mux.Handle("DELETE /api/v1/contacts/{id}",
		middleware.HasSlugs(
			handlers.HandleDeleteContact(
				config.ContactStore, 
				getenv,
			),
			[]string{"id"},
		),
 )

	/* ----- ContactDetails ----- */
	// Read
	mux.Handle("GET /api/v1/contactDetails/{id}",
		middleware.HasSlugs(
			handlers.HandleReadContactDetail(
				config.ContactDetailStore, 
				getenv,
			),
			[]string{"id"},
		),
	)

	/* ----- ContactLists ----- */
	// Read
	mux.Handle("GET /api/v1/contactLists/{id}",
		middleware.HasSlugs(
			handlers.HandleReadContactList(
				config.ContactListStore, 
				getenv,
			),
			[]string{"id"},
		),
	)

	// Write
	mux.Handle("POST /api/v1/contactLists",
		handlers.HandleWriteContactList(
			config.ContactListStore,
			getenv,
		),
	)

	// Update
	mux.Handle("POST /api/v1/contactLists/{id}",
		middleware.HasSlugs(
			handlers.HandleWriteContactListById(
				config.ContactListStore, 
				getenv,
			),
			[]string{"id"},
		),
	)

	// Delete	
	mux.Handle("DELETE /api/v1/contactLists/{id}",
		middleware.HasSlugs(
			handlers.HandleDeleteContactList(
				config.ContactListStore, 
				getenv,
			),
			[]string{"id"},
		),
	)

	// Add contacts
	mux.Handle("POST /api/v1/contactLists/{id}/addContacts",
		middleware.HasSlugs(
			handlers.HandleAddContactsToContactList(
				config.ContactListStore, 
				config.ContactStore,
				getenv,
			),
			[]string{"id"},
		),
	)

	// Remove contacts
	mux.Handle("DELETE /api/v1/contactLists/{id}/removeContacts",
		middleware.HasSlugs(
			handlers.HandleRemoveContactsFromContactList(
				config.ContactListStore, 
				getenv,
			),
			[]string{"id"},
		),
	)
}
