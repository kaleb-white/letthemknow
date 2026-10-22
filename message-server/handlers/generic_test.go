package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kaleb-white/letthemknow/message-server/handlers"
	"github.com/kaleb-white/letthemknow/message-server/mocks"
	"github.com/kaleb-white/letthemknow/message-server/models"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)

// Credit to https://www.cloudbees.com/blog/testing-http-handlers-go

var genericTestMux http.ServeMux = *http.NewServeMux()

func TestGenericReadHandlerGoldenPath(t *testing.T) {
	cs := mocks.ContactStoreThatHasId1AndId2{}
	h := handlers.GenericReadHandler(cs, utils.Getenv, "a", "b")
	
	path := "/TestGenericReadHandlerGoldenPath/1"

	genericTestMux.Handle("/TestGenericReadHandlerGoldenPath/{id}", h)

	r := utils.GenerateRequest(t, "GET", path, nil)
	rr := httptest.NewRecorder()

	genericTestMux.ServeHTTP(rr, r)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler did not succeed, got %v wanted %v", status, http.StatusOK)
		return
	}
	
	receivedContact := models.Contact{}	

	err := json.NewDecoder(rr.Result().Body).Decode(&receivedContact)
	if err != nil {
		t.Errorf("Binary -> struct failed, %s", err.Error())
		return
	}

	actualContact := cs.GetStoredContact()

	errs, ok := actualContact.Equals(&receivedContact)
	if !ok {
		t.Errorf("Did not receive expected contact, errs: %s", utils.PrettyPrintErrors(errs))
	}
}

func TestGenericReadHandlerDbTimeout(t *testing.T)     {
	cs := mocks.StoreThatAlwaysTimesOut{}
	getenv := func(string) string {
		return "100"
	}
	h := handlers.GenericReadHandler(cs, getenv, "a", "b")
	
	path := "/TestGenericReadHandlerDbTimeout/1"

	genericTestMux.Handle("/TestGenericReadHandlerDbTimeout/{id}", h)

	r := utils.GenerateRequest(t, "GET", path, nil)
	rr := httptest.NewRecorder()

	genericTestMux.ServeHTTP(rr, r)
	
	if status := rr.Code; status != mocks.TimeOutStoreTimeOutStatus {
		t.Errorf("Handler did not return expected mock error, got %v wanted %v", status, mocks.TimeOutStoreTimeOutStatus)
		return
	}
	
	receivedErr := utils.CreateHttpError("", http.StatusOK)

	err := json.NewDecoder(rr.Result().Body).Decode(&receivedErr)
	if err != nil {
		t.Errorf("Binary -> struct failed, %s", err.Error())
		return
	}

	if receivedErr.Status != mocks.TimeOutStoreTimeOutStatus {
		t.Errorf("Returned struct did not return expected mock error, got %v wanted %v", receivedErr.Status, mocks.TimeOutStoreTimeOutStatus)
		return
	}
		
}

func TestGenericDeleteHandlerGoldenPath(t *testing.T)  {
	cs := mocks.StoreThatRecordsCalls{}
	h := handlers.GenericDeleteHandler[uint64](cs, utils.Getenv, "a", "b")
	
	path := "/TestGenericDeleteHandlerGoldenPath/1"

	genericTestMux.Handle("/TestGenericDeleteHandlerGoldenPath/{id}", h)

	r := utils.GenerateRequest(t, "GET", path, nil)
	rr := httptest.NewRecorder()

	genericTestMux.ServeHTTP(rr, r)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler did not succeed, got %v wanted %v", status, http.StatusOK)
		return
	}
	
	var receivedContact uint64 

	err := json.NewDecoder(rr.Result().Body).Decode(&receivedContact)
	if err != nil {
		t.Errorf("Binary -> uint64 failed, %s", err.Error())
		return
	}


	if cs.GetNumTimesDeleteCalled() != 1 {
		t.Errorf("Delete was not called right num times, got %v wanted %v", cs.GetNumTimesDeleteCalled(), 1)
	}
}

func TestWriteHandlerGoldenPath(t *testing.T)          {
	cs := mocks.StoreThatRecordsCalls{}
	h := handlers.GenericWriteHandler(cs, utils.Getenv, "WriteHandlerGoldenPath", "contact", func(*models.Contact, *http.Request) {})
	
	path := "/TestGenericWriteHandlerGoldenPath/1"

	genericTestMux.Handle("/TestGenericWriteHandlerGoldenPath/{id}", h)

	r := utils.GenerateRequest(t, "POST", path, models.Contact{
		Id: 2,
		Phone: 66666666,
		FirstName: "asdf",
		LastName: "fdsa",
		CreatedBy: "alalal",
		LastUpdatedBy: "lalalal",
	})
	rr := httptest.NewRecorder()

	genericTestMux.ServeHTTP(rr, r)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler did not succeed, got %v wanted %v", status, http.StatusOK)
		return
	}
	
	var receivedContact models.IdResponse 

	err := json.NewDecoder(rr.Result().Body).Decode(&receivedContact)
	if err != nil {
		t.Errorf("Binary -> models.IdResponse failed, %s", err.Error())
		return
	}

	if cs.GetNumTimesWriteCalled() != 1 {
		t.Errorf("Write was not called right num times, got %v wanted %v", cs.GetNumTimesWriteCalled(), 1)
	}
}

func TestWriteHandlerByIdGoldenPath(t *testing.T)      {
	cs := mocks.ContactStoreThatWillWriteById{}

	validationChange := "sdlkflsdkfjslkdfj"
	h := handlers.GenericWriteByIdHandler(cs, utils.Getenv, "GenericWriteHandlerByIdGoldenPath", "b", func(a *models.Contact, b *http.Request) {a.FirstName = validationChange})
	
	path := "/TestGenericWriteHandlerByIdGoldenPath/1"

	genericTestMux.Handle("/TestGenericWriteHandlerByIdGoldenPath/{id}", h)

	newContact := models.Contact{
		Id: 2,
		Phone: 66666666,
		FirstName: "asdf",
		LastName: "fdsa",
		CreatedBy: "alalal",
		LastUpdatedBy: "lalalal",
	}

	r := utils.GenerateRequest(t, "POST", path, newContact)
	rr := httptest.NewRecorder()

	// Write contact into request
	rr.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rr).Encode(newContact); err != nil {
		t.Errorf("Error during encode: %s", err.Error())
	}

	genericTestMux.ServeHTTP(rr, r)
	
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler did not succeed, got %v wanted %v", status, http.StatusOK)
		return
	}

	receivedContact := cs.GetStoredContact()

	if receivedContact.FirstName != validationChange {
		t.Errorf("Expected first name to change in prevalidation run, got %s wanted %s", receivedContact.FirstName, validationChange)
	}

	newContact.FirstName = validationChange
	errs, ok := receivedContact.Equals(&newContact)
	if !ok {
		t.Errorf("Failed to write contact, errs: %s", utils.PrettyPrintErrors(errs))
	}

}
func TestWriteHandlerByIdMissingId(t *testing.T)			 {
	cs := mocks.ContactStoreThatWillWriteById{}

	h := handlers.GenericWriteHandler(cs, utils.Getenv, "a", "b", func(*models.Contact, *http.Request) {})
	
	path := "/TestGenericWriteHandlerByIdMissingId/45"

	genericTestMux.Handle("/TestGenericWriteHandlerByIdMissingId/{id}", h)

	r := utils.GenerateRequest(t, "POST", path, models.Contact{
		Id: 2,
		Phone: 66666666,
		FirstName: "asdf",
		LastName: "fdsa",
		CreatedBy: "alalal",
		LastUpdatedBy: "lalalal",
	})
	rr := httptest.NewRecorder()

	// Set header 
	rr.Header().Set("Content-Type", "application/json")

	// Record contact before we request and after
	b4 := cs.GetStoredContact()

	genericTestMux.ServeHTTP(rr, r)
	
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler did not succeed, got %v wanted %v", status, http.StatusNotFound)
		return
	}

	after := cs.GetStoredContact()

	errs, ok := b4.Equals(&after)
	if !ok {
		t.Errorf("Contact was modified in flight: %s", utils.PrettyPrintErrors(errs))
	}
}
