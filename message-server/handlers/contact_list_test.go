package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/kaleb-white/letthemknow/message-server/handlers"
	"github.com/kaleb-white/letthemknow/message-server/mocks"
	"github.com/kaleb-white/letthemknow/message-server/models"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)

var contactListTestMux http.ServeMux = *http.NewServeMux()

func TestContactListAddIdsGoldenPath(t *testing.T) {
	cs := mocks.ContactStoreThatHasId1AndId2{}
	cls := mocks.ContactListStoreThatWillWriteById{}
	h := handlers.HandleAddContactsToContactList(cls, cs, utils.Getenv)

	path := "/TestContactListAddIdsGoldenPath/1"

	contactListTestMux.Handle("/TestContactListAddIdsGoldenPath/{id}", h)

	r := utils.GenerateRequest(t, "POST", path, models.AddContactsToListRequest{
		Ids: []uint64{2},
	})
	rr := httptest.NewRecorder()

	clBeforeServe := cls.GetStoredContact()

	contactListTestMux.ServeHTTP(rr, r)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler did not succeed, got %v wanted %v", status, http.StatusOK)
		return
	}

	addedIdsExpected := models.AddedIdsResponse{
		Added: []uint64{2},
	}
	addedIdsActual := models.AddedIdsResponse{}
	err := json.NewDecoder(rr.Result().Body).Decode(&addedIdsActual)
	if err != nil {
		t.Errorf("Binary -> struct failed, %s", err.Error())
		return
	}

	if !slices.Equal(addedIdsActual.Added, addedIdsExpected.Added) {
		if len(addedIdsActual.Added) > 0 {
			t.Errorf("Expected added was not actual, expected '2', got %d", addedIdsActual.Added[0])		
		} else {
			t.Errorf("Expected added was not actual, expected '2', got nothing")
		}
	}
	
	clAfterServe := cls.GetStoredContact()
	clBeforeServe.ContactList = []uint64{1, 2}

	if errs, ok := clAfterServe.Equals(&clBeforeServe); !ok {
		expected := utils.IdSliceToString(clBeforeServe.ContactList)
		actual := utils.IdSliceToString(clAfterServe.ContactList)
		t.Errorf("Contact lists were not equal with expected change after serve, %s. Got %s, wanted %s", utils.PrettyPrintErrors(errs), actual, expected)
	}
}

func TestContactListRemoveIdsGoldenPath(t *testing.T) {
	cls := mocks.ContactListStoreThatWillWriteById{}
	h := handlers.HandleRemoveContactsFromContactList(cls, utils.Getenv)

	path := "/TestContactListRemoveIdsGoldenPath/1"

	contactListTestMux.Handle("/TestContactListRemoveIdsGoldenPath/{id}", h)

	r := utils.GenerateRequest(t, "POST", path, models.RemoveContactsFromListRequest{
		Ids: []uint64{1, 2},
	})
	rr := httptest.NewRecorder()

	clBeforeServe := cls.GetStoredContact()

	contactListTestMux.ServeHTTP(rr, r)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler did not succeed, got %v wanted %v", status, http.StatusOK)
		return
	}

	removedIdsExpected := models.RemovedIdsResponse{
		Removed: []uint64{1, 2},
	}
	removedIdsActual := models.RemovedIdsResponse{}
	err := json.NewDecoder(rr.Result().Body).Decode(&removedIdsActual)
	if err != nil {
		t.Errorf("Binary -> struct failed, %s", err.Error())
		return
	}

	if !slices.Equal(removedIdsActual.Removed, removedIdsExpected.Removed) {
		if len(removedIdsActual.Removed) > 0 {
			t.Errorf("Expected removed was not actual, expected '', got %d", removedIdsActual.Removed[0])		
			return
		}
	}
	
	clAfterServe := cls.GetStoredContact()
	clBeforeServe.ContactList = []uint64{}

	if errs, ok := clAfterServe.Equals(&clBeforeServe); !ok {
		expected := utils.IdSliceToString(clBeforeServe.ContactList)
		actual := utils.IdSliceToString(clAfterServe.ContactList)
		t.Errorf("Contact lists were not equal with expected change after serve, %s. Got %s, wanted %s", utils.PrettyPrintErrors(errs), actual, expected)
		return
	}
}

func TestContactListRemoveIdsBadRequest(t *testing.T) {
	cls := mocks.ContactListStoreThatWillWriteById{}
	h := handlers.HandleRemoveContactsFromContactList(cls, utils.Getenv)

	path := "/TestContactListRemoveIdsBadRequest/1"

	contactListTestMux.Handle("/TestContactListRemoveIdsBadRequest/{id}", h)

	r := utils.GenerateRequest(t, "POST", path, models.RemoveContactsFromListRequest{
		Ids: []uint64{},
	})
	rr := httptest.NewRecorder()

	clBeforeServe := cls.GetStoredContact()

	contactListTestMux.ServeHTTP(rr, r)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler did not succeed, got %v wanted %v", status, http.StatusBadRequest)
		return
	}

	clAfterServer := cls.GetStoredContact()

	if errs, ok := clBeforeServe.Equals(&clAfterServer); !ok {
		t.Errorf("Contact was modified with a bad write, %s", utils.PrettyPrintErrors(errs))
	}
}

func TestContactListAddMultipleIds(t *testing.T) {
	cs := mocks.ContactStoreThatHasId1AndId2{}
	cls := mocks.ContactListStoreThatWillWriteById{}
	h := handlers.HandleAddContactsToContactList(cls, cs, utils.Getenv)

	path := "/TestContactListAddMultipleIds/1"

	contactListTestMux.Handle("/TestContactListAddMultipleIds/{id}", h)

	r := utils.GenerateRequest(t, "POST", path, models.AddContactsToListRequest{
		Ids: []uint64{1, 2},
	})
	rr := httptest.NewRecorder()

	clBeforeServe := cls.GetStoredContact()

	contactListTestMux.ServeHTTP(rr, r)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler did not succeed, got %v wanted %v", status, http.StatusOK)
		return
	}

	addedIdsExpected := models.AddedIdsResponse{
		Added: []uint64{1, 2},
	}
	addedIdsActual := models.AddedIdsResponse{}
	err := json.NewDecoder(rr.Result().Body).Decode(&addedIdsActual)
	if err != nil {
		t.Errorf("Binary -> struct failed, %s", err.Error())
		return
	}

	if !slices.Equal(addedIdsActual.Added, addedIdsExpected.Added) {
		t.Errorf("Expected added was not actual, expected %s, got %s", utils.IdSliceToString(addedIdsExpected.Added), utils.IdSliceToString(addedIdsActual.Added))		
		return
	}
	
	clAfterServe := cls.GetStoredContact()
	clBeforeServe.ContactList = []uint64{1, 2}

	if errs, ok := clAfterServe.Equals(&clBeforeServe); !ok {
		expected := utils.IdSliceToString(clBeforeServe.ContactList)
		actual := utils.IdSliceToString(clAfterServe.ContactList)
		t.Errorf("Contact lists were not equal with expected change after serve, %s. Got %s, wanted %s", utils.PrettyPrintErrors(errs), actual, expected)
	}
}

func TestContactListAddNonExistentId(t *testing.T) {
	cs := mocks.ContactStoreThatHasId1AndId2{}
	cls := mocks.ContactListStoreThatWillWriteById{}
	h := handlers.HandleAddContactsToContactList(cls, cs, utils.Getenv)

	path := "/TestContactListAddNonExistentId/1"

	contactListTestMux.Handle("/TestContactListAddNonExistentId/{id}", h)

	r := utils.GenerateRequest(t, "POST", path, models.AddContactsToListRequest{
		Ids: []uint64{3},
	})
	rr := httptest.NewRecorder()

	contactListTestMux.ServeHTTP(rr, r)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler did not succeed, got %v wanted %v", status, http.StatusBadRequest)
		return
	}
}
