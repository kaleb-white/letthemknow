package cd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/kaleb-white/letthemknow/message-server/db"
	"github.com/kaleb-white/letthemknow/message-server/log"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)

const SOURCE_CONTACT_DETAIL = "contact_detail"

type ContactDetail struct {
	Id uint64
	Phone uint64
	FirstName string
	LastName string
	LastUpdatedAtDateTime time.Time
	LastUpdatedBy string
	Org string	
}

func (c *ContactDetail) Validate() (bool, []error) {
	collectedErrors := make([]error, 0, 10)

	// Check requiredFields
	requiredFields := []string{"Phone", "FirstName", "LastName", "LastUpdatedBy"}
	wasError := utils.CheckRequiredFieldsArentDefault(c, &requiredFields, &collectedErrors)

	if c.LastUpdatedAtDateTime.Compare(time.Now()) == 1 {
		collectedErrors = append(collectedErrors, errors.New("LastUpdatedAtDateTime must be in the past (use time.Now())."))
		wasError = true
	}

	return wasError, collectedErrors
}

var CONTACT_DETAIL_PREPARED_SQL = map[string]db.PreparedSql{
	"Read": {
		OperationName: "contact_details_read",
		RawSql: CONTACT_DETAIL_READ,
	},
	"Delete": {
		OperationName: "contact_details_delete",
		RawSql: CONTACT_DETAIL_DELETE,
	},
}

func (c *ContactDetail) Read(ctx context.Context, Id uint64) (*ContactDetail, error) {
	// Connect
	conn, err := db.GetConnection()
	if err != nil {
		log.Log(SOURCE_CONTACT_DETAIL, fmt.Sprintf("Failed to connect to read contact details with id %d", Id), log.INFO)
		return c, err
	}
	defer conn.Close()

	// Use statement singleton
	preparedSql := CONTACT_DETAIL_PREPARED_SQL["Read"]
	stmt, err := preparedSql.GetPreparedStatement(ctx, conn)
	if err != nil {
		log.Log(SOURCE_CONTACT_DETAIL, fmt.Sprintf("Failed to prepare stmt to read contact details for id %d: %s", Id, err.Error()), log.ERROR)
		return c, err
	}

	// Query 
	var LastUpdatedAtDateTime []byte 
	err = stmt.QueryRowContext(ctx, Id).Scan(&c.Id, &c.Phone, &c.FirstName, &c.LastName, &LastUpdatedAtDateTime, &c.LastUpdatedBy, &c.Org)

	if err == sql.ErrNoRows {
		log.Log(SOURCE_CONTACT_DETAIL, fmt.Sprintf("No rows found while reading id %d", Id), log.WARNING)
		return c, err
	} else if err != nil {
		log.Log(SOURCE_CONTACT_DETAIL, fmt.Sprintf("Failed to execute stmt to read contact details for id %d: %s", Id, err.Error()), log.ERROR)
		return c, err
	}

	// Map
	err = c.LastUpdatedAtDateTime.UnmarshalBinary(LastUpdatedAtDateTime)
	if err != nil {
		log.Log(SOURCE_CONTACT_DETAIL, fmt.Sprintf("Failed to unmarshal LastUpdatedAtDateTime. Raw: %s", LastUpdatedAtDateTime), log.WARNING)
	}

	return c, nil
}

func (c *ContactDetail) Delete(ctx context.Context) (uint64, error) {
	// Connect
	conn, err := db.GetConnection()
	if err != nil {
		log.Log(SOURCE_CONTACT_DETAIL, fmt.Sprintf("Failed to connect to delete contact details with id %d", c.Id), log.INFO)
		return 0, err
	}
	defer conn.Close()

	// Use statement singleton
	preparedSql := CONTACT_DETAIL_PREPARED_SQL["Delete"]
	stmt, err := preparedSql.GetPreparedStatement(ctx, conn)
	if err != nil {
		log.Log(SOURCE_CONTACT_DETAIL, fmt.Sprintf("Failed to prepare stmt to delete contact details for id %d: %s", c.Id, err.Error()), log.ERROR)
		return 0, err
	}

	// Query 
	var Id uint64 
	err = stmt.QueryRowContext(ctx, &c.Id).Scan(&Id)
	if err == sql.ErrNoRows {
		log.Log(SOURCE_CONTACT_DETAIL, fmt.Sprintf("No rows found while reading id %d", Id), log.WARNING)
		return 0, err
	} else if c.Id != Id {
		log.Log(SOURCE_CONTACT_DETAIL, fmt.Sprintf("Deleted id did not match what was passed. Got: %d, expected: %d.", Id, c.Id), log.ERROR)
		return 0, err
	} else if err != nil {
		log.Log(SOURCE_CONTACT_DETAIL, fmt.Sprintf("Failed to execute stmt to read contact for id %d: %s", Id, err.Error()), log.ERROR)
		return 0, err
	}

	return Id, nil
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
