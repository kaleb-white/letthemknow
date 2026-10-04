package models

import (
	"time"
	"context"
	"errors"
	"fmt"

	"github.com/kaleb-white/letthemknow/message-server/utils"
)

type ContactDetail struct {
	Id uint64
	Phone uint64
	FirstName string
	LastName string
	LastUpdatedAtDateTime time.Time
	LastUpdatedBy string
	Org string	
}

type ContactDetailStore interface {
	Init(context.Context) error
	Read(context.Context, uint64) (ContactDetail, utils.StoreError)
	Delete(context.Context, uint64) (uint64, utils.StoreError)
}

func (c *ContactDetail) SetId(Id uint64) {
	c.Id = Id
}

func (c *ContactDetail) GetId() uint64 {
	return c.Id
}

func (c *ContactDetail) Validate() (bool, []error) {
	collectedErrors := make([]error, 0, 10)

	// Check requiredFields
	requiredFields := []string{"Phone", "FirstName", "LastName", "LastUpdatedBy"}
	wasError := utils.CheckRequiredFieldsArentDefault(*c, &requiredFields, &collectedErrors)

	if c.LastUpdatedAtDateTime.Compare(time.Now()) == 1 {
		collectedErrors = append(collectedErrors, errors.New("LastUpdatedAtDateTime must be in the past (use time.Now())."))
		wasError = true
	}

	return wasError, collectedErrors
}

func (c *ContactDetail) Equals(c2 *ContactDetail) ([]error, bool) {
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
	if c2.FirstName != c.FirstName {
		wasError = true
		errs = append(errs, fmt.Errorf("FirstName: %s, expected %s", c2.FirstName, c.FirstName))
	}
	if c2.LastName != c.LastName {
		wasError = true
		errs = append(errs, fmt.Errorf("LastName: %s, expected %s", c2.LastName, c.LastName))
	}
	if c2.Org != c.Org {
		wasError = true
		errs = append(errs, errors.New("Org"))
	}
	return errs, !wasError
}
