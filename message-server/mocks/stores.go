package mocks

import (
	"context"
	"net/http"
	"time"

	"github.com/kaleb-white/letthemknow/message-server/models"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)

type StoreThatAlwaysTimesOut struct {
}

func (s StoreThatAlwaysTimesOut) Init(c context.Context) error {
	deadline, hasDeadline := c.Deadline()
	if !hasDeadline {
		return nil
	}
	
	go func(){
		time.Sleep(time.Until(deadline) + time.Second)
	}()

	<- c.Done()
	return c.Err()
}

// Use uncommon status to ensure tests source from this const
const TimeOutStoreTimeOutStatus int = http.StatusConflict

func (s StoreThatAlwaysTimesOut) Write(c context.Context, _ *int) (uint64, utils.StoreError) {
	deadline, hasDeadline := c.Deadline()
	if !hasDeadline {
		return 0, nil
	}
	
	go func(){
		time.Sleep(time.Until(deadline) + time.Second)
	}()

	<- c.Done()
	return 0, utils.NewStoreError(c.Err().Error(), TimeOutStoreTimeOutStatus)
}

func (s StoreThatAlwaysTimesOut) Read(c context.Context, _ uint64) (int, utils.StoreError) {
	deadline, hasDeadline := c.Deadline()
	if !hasDeadline {
		return 0, nil
	}
	
	go func(){
		time.Sleep(time.Until(deadline) + time.Second)
	}()

	<- c.Done()
	return 0, utils.NewStoreError(c.Err().Error(), TimeOutStoreTimeOutStatus)
}

func (s StoreThatAlwaysTimesOut) Delete(c context.Context, _ uint64) (uint64, utils.StoreError) {
	deadline, hasDeadline := c.Deadline()
	if !hasDeadline {
		return 0, nil
	}
	
	go func(){
		time.Sleep(time.Until(deadline) + time.Second)
	}()

	<- c.Done()
	return 0, utils.NewStoreError(c.Err().Error(), TimeOutStoreTimeOutStatus)
}

func (_ StoreThatAlwaysTimesOut) ReadTimeout(func (string) string) time.Duration {
	return time.Duration(1000)
}
func (_ StoreThatAlwaysTimesOut) WriteTimeout(func (string) string) time.Duration {
	return time.Duration(1000)
}
func (_ StoreThatAlwaysTimesOut) DeleteTimeout(func (string) string) time.Duration {
	return time.Duration(1000)
}


////////////////////////////////////////////////////////////////////////
type StoreThatReturnsZero struct {
}

func (s StoreThatReturnsZero) Init(c context.Context) error {
	return nil
}

func (s StoreThatReturnsZero) Write(c context.Context, _ *int) (uint64, utils.StoreError) {
	return 0, nil
}

func (s StoreThatReturnsZero) Read(c context.Context, _ uint64) (int, utils.StoreError) {
	return 0, nil
}

func (s StoreThatReturnsZero) Delete(c context.Context, _ uint64) (uint64, utils.StoreError) {
	return 0, nil
}

func (_ StoreThatReturnsZero) ReadTimeout(func (string) string) time.Duration {
	return time.Duration(1000)
}
func (_ StoreThatReturnsZero) WriteTimeout(func (string) string) time.Duration {
	return time.Duration(1000)
}
func (_ StoreThatReturnsZero) DeleteTimeout(func (string) string) time.Duration {
	return time.Duration(1000)
}


////////////////////////////////////////////////////////////////////////
type ContactStoreThatReturnsZero struct {
}

func (s ContactStoreThatReturnsZero) Init(c context.Context) error {
	return nil
}

func (s ContactStoreThatReturnsZero) Write(c context.Context, _ *models.Contact) (uint64, utils.StoreError) {
	return 0, nil
}

func (s ContactStoreThatReturnsZero) Read(c context.Context, _ uint64) (models.Contact, utils.StoreError) {
	return models.Contact{
		Id: 1,
		Phone: 1234567891,
		FirstName: "big",
		LastName: "boy",
		CreatedBy: "teseter",
		LastUpdatedBy: "test",
	}, nil
}

func (s ContactStoreThatReturnsZero) Delete(c context.Context, _ uint64) (uint64, utils.StoreError) {
	return 0, nil
}

func (_ ContactStoreThatReturnsZero) ReadTimeout(func (string) string) time.Duration {
	return time.Duration(1000)
}
func (_ ContactStoreThatReturnsZero) WriteTimeout(func (string) string) time.Duration {
	return time.Duration(1000)
}
func (_ ContactStoreThatReturnsZero) DeleteTimeout(func (string) string) time.Duration {
	return time.Duration(1000)
}


////////////////////////////////////////////////////////////////////////
type ContactStoreThatHasId1AndId2 struct {
}

var contactWithId1 models.Contact = models.Contact{
	Id: 1,
	Phone: 1234567891,
	FirstName: "big",
	LastName: "boy",
	CreatedBy: "teseter",
	LastUpdatedBy: "test",
}

var contactWithId2 models.Contact = models.Contact{
	Id: 2,
	Phone: 1234567891,
	FirstName: "big",
	LastName: "man",
	CreatedBy: "teseter",
	LastUpdatedBy: "test",

}

func (s ContactStoreThatHasId1AndId2) Init(c context.Context) error {
	return nil
}

func (s ContactStoreThatHasId1AndId2) Write(c context.Context, _ *models.Contact) (uint64, utils.StoreError) {
	return 0, nil
}

func (s ContactStoreThatHasId1AndId2) Read(c context.Context, id uint64) (models.Contact, utils.StoreError) {
	switch id {
	case 1:
		return contactWithId1, nil
	case 2:
		return contactWithId2, nil
	}
	
	return models.Contact{}, utils.NewStoreError("not found", http.StatusNotFound)
}

func (s ContactStoreThatHasId1AndId2) Delete(c context.Context, _ uint64) (uint64, utils.StoreError) {
	return 0, nil
}

func (s ContactStoreThatHasId1AndId2) GetStoredContact() models.Contact {
	return contactWithId1
}
func (_ ContactStoreThatHasId1AndId2) ReadTimeout(func (string) string) time.Duration {
	return time.Duration(1000)
}
func (_ ContactStoreThatHasId1AndId2) WriteTimeout(func (string) string) time.Duration {
	return time.Duration(1000)
}
func (_ ContactStoreThatHasId1AndId2) DeleteTimeout(func (string) string) time.Duration {
	return time.Duration(1000)
}


////////////////////////////////////////////////////////////////////////
type StoreThatRecordsCalls struct {
}

func (s StoreThatRecordsCalls) Init(c context.Context) error {
	return nil
}

var timesWriteCalled int = 0

func (s StoreThatRecordsCalls) GetNumTimesWriteCalled() int {
	return timesWriteCalled
}

func (s StoreThatRecordsCalls) Write(c context.Context, _ *models.Contact) (uint64, utils.StoreError) {
	timesWriteCalled++
	return 0, nil
}

var timesReadCalled int = 0

func (s StoreThatRecordsCalls) GetNumTimesReadCalled() int {
	return timesReadCalled
}

func (s StoreThatRecordsCalls) Read(c context.Context, _ uint64) (models.Contact, utils.StoreError) {
	timesReadCalled++
	return models.Contact{}, nil
}

var timesDeleteCalled int = 0

func (s StoreThatRecordsCalls) GetNumTimesDeleteCalled() int {
	return timesDeleteCalled
}

func (s StoreThatRecordsCalls) Delete(c context.Context, _ uint64) (uint64, utils.StoreError) {
	timesDeleteCalled++
	return 0, nil
}
func (_ StoreThatRecordsCalls) ReadTimeout(func (string) string) time.Duration {
	return time.Duration(1000)
}
func (_ StoreThatRecordsCalls) WriteTimeout(func (string) string) time.Duration {
	return time.Duration(1000)
}
func (_ StoreThatRecordsCalls) DeleteTimeout(func (string) string) time.Duration {
	return time.Duration(1000)
}


////////////////////////////////////////////////////////////////////////
type ContactStoreThatWillWriteById struct {
}

var writeableContactWithId1 models.Contact = models.Contact{
	Id: 1,
	Phone: 1234567891,
	FirstName: "big",
	LastName: "boy",
	CreatedBy: "teseter",
	LastUpdatedBy: "test",
}

func (s ContactStoreThatWillWriteById) Init(c context.Context) error {
	return nil
}

func (s ContactStoreThatWillWriteById) Write(c context.Context, new *models.Contact) (uint64, utils.StoreError) {
	if new.Id == writeableContactWithId1.Id {
		writeableContactWithId1 = *new
		return writeableContactWithId1.Id, nil
	} else {
		return 0, utils.NewStoreError("not found", http.StatusNotFound)
	}
}

func (s ContactStoreThatWillWriteById) GetStoredContact() models.Contact {
	return writeableContactWithId1
}

func (_ ContactStoreThatWillWriteById) ReadTimeout(func (string) string) time.Duration {
	return time.Duration(1000)
}
func (_ ContactStoreThatWillWriteById) WriteTimeout(func (string) string) time.Duration {
	return time.Duration(1000)
}
func (_ ContactStoreThatWillWriteById) DeleteTimeout(func (string) string) time.Duration {
	return time.Duration(1000)
}

////////////////////////////////////////////////////////////////////////
type ContactListStoreThatWillWriteById struct {
	cl models.ContactList
	initialized bool
}

var writeableContactListWithId1 models.ContactList = models.ContactList{
	Id: 1,
	Name:  "MyList",	
	ContactList: []uint64{1},
	CreatedBy: "teseter",
	LastUpdatedBy: "test",
}

func (s ContactListStoreThatWillWriteById) Init(c context.Context) error {
	return nil
}

func (s ContactListStoreThatWillWriteById) Write(c context.Context, new *models.ContactList) (uint64, utils.StoreError) {
	if new.Id == writeableContactListWithId1.Id {
		writeableContactListWithId1 = *new
		return writeableContactListWithId1.Id, nil
	} else {
		return 0, utils.NewStoreError("not found", http.StatusNotFound)
	}
}

func (s ContactListStoreThatWillWriteById) Read(context.Context, uint64) (models.ContactList, utils.StoreError) {
	return writeableContactListWithId1, nil
}

func (s ContactListStoreThatWillWriteById) Delete(context.Context, uint64) (uint64, utils.StoreError) {
	id := writeableContactListWithId1.Id
	writeableContactListWithId1 = models.ContactList{}
	return id, nil
}

func (s ContactListStoreThatWillWriteById) GetStoredContact() models.ContactList {
	return writeableContactListWithId1
}

func (_ ContactListStoreThatWillWriteById) ReadTimeout(func (string) string) time.Duration {
	return time.Duration(1000)
}
func (_ ContactListStoreThatWillWriteById) WriteTimeout(func (string) string) time.Duration {
	return time.Duration(1000)
}
func (_ ContactListStoreThatWillWriteById) DeleteTimeout(func (string) string) time.Duration {
	return time.Duration(1000)
}
