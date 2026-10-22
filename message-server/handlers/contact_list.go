package handlers

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/kaleb-white/letthemknow/message-server/log"
	"github.com/kaleb-white/letthemknow/message-server/models"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)

const SOURCE_CONTACT_LIST_HANDLER string = "contact_list_handler"

func HandleReadContactList(contactListStore models.ContactListStore, getenv func(string) string) http.Handler {
	return GenericReadHandler(contactListStore, getenv, SOURCE_CONTACT_LIST_HANDLER, "ContactList")
}

func HandleWriteContactList(contactListStore models.ContactListStore, getenv func(string) string) http.Handler {
	return GenericWriteHandler(
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

func HandleWriteContactListById(contactListStore models.ContactListStore, getenv func(string) string) http.Handler {
	return GenericWriteByIdHandler(
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

func HandleDeleteContactList(contactListStore models.ContactListStore, getenv func(string) string) http.Handler {
	return GenericDeleteHandler[models.Contact](contactListStore, getenv, SOURCE_CONTACT_HANDLERS, "Contact")
}

func HandleAddContactsToContactList(
	contactListStore models.ContactListStore, 
	contactStore models.ContactStore, 
	getenv func(string) string,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// This should already be checked by HasSlugs middleware
			id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)

			// Perform request to store
			contactList, readErr := ReadFromStoreWithContext(id, contactListStore, getenv)

			// Parse result
			if readErr != nil {
				e := fmt.Sprintf("Failure to read contact list: %s", readErr.Error())
				log.Log(SOURCE_CONTACT_LIST_HANDLER, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(readErr.Error(), readErr.Status(), w, r)
				return
			}

			idList, err := utils.Decode[models.AddContactsToListRequest](r)	

			if err != nil {
				e := fmt.Sprintf("Failure while decoding contacts to add them to list: %s", err.Error())
				log.Log(SOURCE_CONTACT_LIST_HANDLER, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(err.Error(), 400, w, r)
				return
			}

			// Validate list has ids
			if len(idList.Ids) == 0 {
				e := "List must contain more than one id"
				log.Log(SOURCE_CONTACT_LIST_HANDLER, utils.AttributeLog(e, r), log.INFO)	
				utils.EncodeHttpError(e, http.StatusBadRequest, w, r)
				return
			}
			
			// Ensure each contact exists and isn't present already
			readErrs := []error{}
			newIds := []uint64{}
			var wg sync.WaitGroup
			for _, id := range idList.Ids {
				if slices.Contains(contactList.ContactList, id) {
					continue
				}
				
				// Async check existence
				wg.Go(func() {
					c, err := ReadFromStoreWithContext(id, contactStore, getenv)	
					if err != nil {
						readErrs = append(readErrs, err)
					} else {
						newIds = append(newIds, c.Id)
					}
				},)
			}

			// Wait for reads to finish
			wg.Wait()

			// Return on failures
			if len(readErrs) > 0 {
				e := utils.PrettyPrintErrors(readErrs)
				log.Log(SOURCE_CONTACT_LIST_HANDLER, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(e, http.StatusBadRequest, w, r)
				return
			}

			// Save
			contactList.ContactList = append(contactList.ContactList, newIds...)
			Id, storeErr := WriteToStoreWithContext(&contactList, contactListStore, getenv)
			if storeErr != nil {
				e := fmt.Sprintf("Error while saving contact list after add: %s", storeErr.Error())
				log.Log(SOURCE_CONTACT_LIST_HANDLER, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(e, storeErr.Status(), w, r)
				return
			}
			
			resp := models.AddedIdsResponse{
				Added: newIds,
			}		
			encodeErr := utils.Encode(w, r, http.StatusOK, resp)
			if encodeErr != nil {
				e := "Failed to encode return object"
				log.Log(SOURCE_CONTACT_LIST_HANDLER, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(e, http.StatusInternalServerError, w, r)
			}

			log.Log(SOURCE_CONTACT_LIST_HANDLER, utils.AttributeLog(fmt.Sprintf("Added ids to contact list with id %d", Id), r), log.INFO)
		},
	)
}

func HandleRemoveContactsFromContactList(
	contactListStore models.ContactListStore, 
	getenv func(string) string,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// This should already be checked by HasSlugs middleware
			id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)

			// Perform request to store
			contactList, readErr := ReadFromStoreWithContext(id, contactListStore, getenv)

			// Parse result
			if readErr != nil {
				e := fmt.Sprintf("Failure to read contact list: %s", readErr.Error())
				log.Log(SOURCE_CONTACT_LIST_HANDLER, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(readErr.Error(), readErr.Status(), w, r)
				return
			}

			idList, err := utils.Decode[models.RemoveContactsFromListRequest](r)	

			if err != nil {
				e := fmt.Sprintf("Failure while decoding contacts to remove them from list: %s", err.Error())
				log.Log(SOURCE_CONTACT_LIST_HANDLER, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(err.Error(), 400, w, r)
				return
			}

			// Validate list has ids
			if len(idList.Ids) == 0 {
				e := "List must contain more than one id"
				log.Log(SOURCE_CONTACT_LIST_HANDLER, utils.AttributeLog(e, r), log.INFO)	
				utils.EncodeHttpError(e, http.StatusBadRequest, w, r)
				return
			}
			
			// Remove from contact list
			removed := []uint64{}
			contactList.ContactList = slices.DeleteFunc(
				contactList.ContactList,
				func(id uint64) bool { 
					if slices.Contains(idList.Ids, id) {
						removed = append(removed, id)
						return true
					}
					return false
				},
			)

			// Save
			Id, storeErr := WriteToStoreWithContext(&contactList, contactListStore, getenv)
			if storeErr != nil {
				e := err.Error()
				log.Log(SOURCE_CONTACT_LIST_HANDLER, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(e, storeErr.Status(), w, r)
				return
			}
			
			resp := models.RemovedIdsResponse{
				Removed: removed,
			}
			encodeErr := utils.Encode(w, r, http.StatusOK, resp)
			if encodeErr != nil {
				e := "Failed to encode return object"
				log.Log(SOURCE_CONTACT_LIST_HANDLER, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(e, http.StatusInternalServerError, w, r)
			}

			log.Log(SOURCE_CONTACT_LIST_HANDLER, utils.AttributeLog(fmt.Sprintf("Removed ids from contact list with id %d", Id), r), log.INFO)
		},
	)
}
