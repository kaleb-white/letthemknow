package schemas

import (
	"errors"
	"time"

	"github.com/kaleb-white/letthemknow/utils"
)

type Contact struct {
	Id uint64
	Phone uint64
	Phone2 uint64
	Phone3 uint64
	FirstName string
	LastName string
	FullName string
	CreatedAtDateTime time.Time
	CreatedBy string
	LastUpdatedAtDateTime time.Time
	LastUpdatedBy string
	Org string	
}

func (c *Contact) Validate() (bool, []error) {
	collectedErrors := make([]error, 0, 10)

	// Check requiredFields
	requiredFields := []string{"Phone", "FirstName", "CreatedBy", "LastUpdatedBy"}
	wasError := utils.CheckRequiredFieldsArentDefault(c, &requiredFields, &collectedErrors)

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

func (c *Contact) Read(Id string) (*Contact, error) {}
func (c *Contact) Read(Id string) (*Contact, error) {}
func (c *Contact) Read(Id string) (*Contact, error) {}


