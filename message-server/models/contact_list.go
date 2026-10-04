package models

import (
	"errors"
	"time"
	"context"

	"github.com/kaleb-white/letthemknow/message-server/utils"
)

type ContactList struct {
	Id                    uint64
	Name 									string
	Description						string
	ContactList						[]uint64
	CreatedAtDateTime     time.Time
	CreatedBy             string
	LastUpdatedAtDateTime time.Time
	LastUpdatedBy         string
}

type ContactListStore interface {
	Init(context.Context) error
	Read(context.Context, uint64) (ContactList, utils.StoreError)
	Write(context.Context, *ContactList) (uint64, utils.StoreError)
	Delete(context.Context, uint64) (uint64, utils.StoreError)
}


func (c *ContactList) SetId(Id uint64) {
	c.Id = Id
}

func (c *ContactList) GetId() uint64 {
	return c.Id
}


func (c *ContactList) Validate() (bool, []error) {
	collectedErrors := make([]error, 0, 10)

	// Check required fields
	requiredFields := []string{"Name", "CreatedBy", "LastUpdatedBy"}
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

func (c *ContactList) Equals(c2 *ContactList) ([]error, bool) {
	errs := make([]error, 0, 5)
	wasError := false

	if c2 == nil {
		wasError = true
		errs = append(errs, errors.New("nil pointer"))
		return errs, !wasError
	}

	fields := []string{"Name", "Description", "ContactList", "CreatedBy", "LastUpdatedBy"}
	wasError = utils.CheckFieldEquality(*c, *c2, &fields, &errs)
			
	return errs, !wasError
}
