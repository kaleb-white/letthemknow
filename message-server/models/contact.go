package models

import (
	"context"
	"time"
	"errors"

	"github.com/kaleb-white/letthemknow/message-server/utils"
)

type Contact struct {
	Id uint64
	Phone uint64
	Phone2 uint64
	Phone3 uint64
	FirstName string
	LastName string
	FullName string
	ListMembership []uint64
	CreatedAtDateTime time.Time
	CreatedBy string
	LastUpdatedAtDateTime time.Time
	LastUpdatedBy string
	Org string	
}

type ContactStore interface {
	Init(context.Context) error
	Read(context.Context, uint64) (Contact, utils.StoreError)
	Write(context.Context, *Contact) (uint64, utils.StoreError)
	Delete(context.Context, uint64) (uint64, utils.StoreError)
}

func (c *Contact) SetId(Id uint64) {
	c.Id = Id
}

func (c *Contact) GetId() uint64 {
	return c.Id
}

func (c *Contact) Validate() (bool, []error) {
	collectedErrors := make([]error, 0, 10)

	// Check requiredFields
	requiredFields := []string{"Phone", "FirstName", "CreatedBy", "LastUpdatedBy"}
	wasError := utils.CheckRequiredFieldsArentDefault(*c, &requiredFields, &collectedErrors)

	// Check datetimes are current
	if c.CreatedAtDateTime.Compare(time.Now()) == 1 {
		collectedErrors = append(collectedErrors, errors.New("CreatedAtDateTime must be in the past (use time.Now())."))
		wasError = true
	}

	if c.LastUpdatedAtDateTime.Compare(time.Now()) == 1 {
		collectedErrors = append(collectedErrors, errors.New("LastUpdatedAtDateTime must be in the past (use time.Now())."))
		wasError = true
	}

	return wasError, collectedErrors
}

func (c *Contact) Equals(c2 *Contact) ([]error, bool) {
	errs := make([]error, 0, 11)
	wasError := false

	if c2 == nil {
		wasError = true
		errs = append(errs, errors.New("nil pointer"))
		return errs, !wasError
	}

	if c2.Phone != c.Phone {
		wasError = true
		errs = append(errs, errors.New("Phone"))
	}
	if c2.Phone2 != c.Phone2 {
		wasError = true
		errs = append(errs, errors.New("Phone2"))
	}
	if c2.Phone3 != c.Phone3 {
		wasError = true
		errs = append(errs, errors.New("Phone3"))
	}
	if c2.FirstName != c.FirstName {
		wasError = true
		errs = append(errs, errors.New("FirstName"))
	}
	if c2.LastName != c.LastName {
		wasError = true
		errs = append(errs, errors.New("LastName"))
	}
	if c2.FullName != c.FullName {
		wasError = true
		errs = append(errs, errors.New("FullName"))
	}
	if !c2.CreatedAtDateTime.Equal(c.CreatedAtDateTime) {
		wasError = true
		errs = append(errs, errors.New("CreatedAtDateTime"))
	}
	if c2.CreatedBy != c.CreatedBy {
		wasError = true
		errs = append(errs, errors.New("CreatedBy"))
	}
	if c2.Org != c.Org {
		wasError = true
		errs = append(errs, errors.New("Org"))
	}
	return errs, !wasError
}

