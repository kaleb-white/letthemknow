package cl

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/kaleb-white/letthemknow/message-server/db"
	"github.com/kaleb-white/letthemknow/message-server/log"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)

const SOURCE_CONTACT_LIST string = "contact_list"

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

var CONTACT_LIST_PREPARED_SQL = map[string]db.PreparedSql{
	"Read": {
		OperationName: "contact_list_read",
		RawSql: CONTACT_LIST_READ,
	},
	"WriteNew": {
		OperationName: "contact_list_write",
		RawSql: CONTACT_LIST_WRITE_NEW,
	},
	"WriteExisting": {
		OperationName: "contact_list_write",
		RawSql: CONTACT_LIST_WRITE_EXISTING,
	},
	"Delete": {
		OperationName: "contact_list_delete",
		RawSql: CONTACT_LIST_DELETE,
	},
}

func (c *ContactList) Read(ctx context.Context, Id uint64) (*ContactList, error) {
	// Connect
	conn, err := db.GetConnection()
	if err != nil {
		log.Log(SOURCE_CONTACT_LIST, fmt.Sprintf("Failed to connect to read contact list with id %d", Id), log.INFO)
		return c, err
	}
	defer conn.Close()
	
	// Use statement singleton
	preparedSql := CONTACT_LIST_PREPARED_SQL["Read"]
	stmt, err := preparedSql.GetPreparedStatement(ctx, conn)
	if err != nil {
		log.Log(SOURCE_CONTACT_LIST, fmt.Sprintf("Failed to prepare stmt to read contact list for id %d: %s", Id, err.Error()), log.ERROR)
		return c, err
	}

	// Query 
	var CreatedAtDateTime, LastUpdatedAtDateTime, ContactList []byte 
	err = stmt.QueryRowContext(ctx, Id).Scan(&c.Id, &c.Name, &c.Description, &ContactList, &CreatedAtDateTime, &c.CreatedBy, &LastUpdatedAtDateTime, &c.LastUpdatedBy)
	if err == sql.ErrNoRows {
		log.Log(SOURCE_CONTACT_LIST, fmt.Sprintf("No rows found while reading id %d", Id), log.WARNING)
		return c, err
	} else if err != nil {
		log.Log(SOURCE_CONTACT_LIST, fmt.Sprintf("Failed to execute stmt to read contact list for id %d: %s", Id, err.Error()), log.ERROR)
		return c, err
	}

	// Map
	err =	c.CreatedAtDateTime.UnmarshalBinary(CreatedAtDateTime)
	if err != nil {
		log.Log(SOURCE_CONTACT_LIST, fmt.Sprintf("Failed to unmarshal CreatedAtDateTime. Raw: %s", CreatedAtDateTime), log.WARNING)
	}
	err = c.LastUpdatedAtDateTime.UnmarshalBinary(LastUpdatedAtDateTime)
	if err != nil {
		log.Log(SOURCE_CONTACT_LIST, fmt.Sprintf("Failed to unmarshal LastUpdatedAtDateTime. Raw: %s", LastUpdatedAtDateTime), log.WARNING)
	}
	err = json.Unmarshal(ContactList, &c.ContactList)
	if err != nil {
		log.Log(SOURCE_CONTACT_LIST, fmt.Sprintf("Failed to unmarshal ContactList. Raw: %s. Err: %s", LastUpdatedAtDateTime, err.Error()), log.WARNING)
	}

	return c, nil
}

func (c *ContactList) Write(ctx context.Context) (uint64, error) {
	// Connect
	conn, err := db.GetConnection()
	if err != nil {
		log.Log(SOURCE_CONTACT_LIST, fmt.Sprintf("Failed to connect to write contact list with id %d", c.Id), log.INFO)
		return 0, err
	}
	defer conn.Close()
	
	// Use statement singleton: zero value Id means new
	var preparedSql db.PreparedSql	
	if c.Id == 0 {
		preparedSql = CONTACT_LIST_PREPARED_SQL["WriteNew"]
	} else {
		preparedSql = CONTACT_LIST_PREPARED_SQL["WriteExisting"]
	}
	stmt, err := preparedSql.GetPreparedStatement(ctx, conn)
	if err != nil {
		log.Log(SOURCE_CONTACT_LIST, fmt.Sprintf("Failed to prepare stmt to write contact list for id %d: %s", c.Id, err.Error()), log.ERROR)
		return 0, err
	}

	// Map
	CreatedAtDateTime, err := c.CreatedAtDateTime.MarshalBinary()
	if err != nil {
		log.Log(SOURCE_CONTACT_LIST, fmt.Sprintf("Failed to marshal CreatedAtDateTime. Value of CreatedAtDateTime: %s", c.CreatedAtDateTime), log.ERROR)
		return 0, err
	}
	LastUpdatedAtDateTime, err := time.Now().MarshalBinary()
	if err != nil {
		log.Log(SOURCE_CONTACT_LIST, fmt.Sprintf("Failed to marshal LastUpdatedAtDateTime. Value of LastUpdatedAtDateTime: %s", c.LastUpdatedAtDateTime), log.ERROR)
		return 0, err
	}
  ContactList, err := json.Marshal(c.ContactList)

	// Query
	var Id uint64
	if c.Id == 0 {
		err = stmt.QueryRowContext(ctx, &c.Name, &c.Description, &ContactList, &CreatedAtDateTime, &c.CreatedBy, &LastUpdatedAtDateTime, &c.LastUpdatedBy).Scan(&Id)
	} else {
		err = stmt.QueryRowContext(ctx, &c.Id, &c.Name, &c.Description, &ContactList, &CreatedAtDateTime, &c.CreatedBy, &LastUpdatedAtDateTime, &c.LastUpdatedBy).Scan(&Id)
	}

	if err == sql.ErrNoRows {
		log.Log(SOURCE_CONTACT_LIST, fmt.Sprintf("No rows returned after writing contact list %s", c.Name), log.WARNING)
		return Id, err
	} else if err != nil {
		log.Log(SOURCE_CONTACT_LIST, fmt.Sprintf("Failed to execute stmt to write for contact list %s: %s", c.Name, err.Error()), log.ERROR)
		return Id, err
	}

	c.Id = Id

	return Id, nil
}

func (c *ContactList) Delete(ctx context.Context) (uint64, error) {
	// Connect
	conn, err := db.GetConnection()
	if err != nil {
		log.Log(SOURCE_CONTACT_LIST, fmt.Sprintf("Failed to connect to delete ContactList with id %d", c.Id), log.INFO)
		return 0, err
	}
	defer conn.Close()
	
	// Use statement singleton
	preparedSql := CONTACT_LIST_PREPARED_SQL["Delete"]
	stmt, err := preparedSql.GetPreparedStatement(ctx, conn)
	if err != nil {
		log.Log(SOURCE_CONTACT_LIST, fmt.Sprintf("Failed to prepare stmt to delete contact list for id %d: %s", c.Id, err.Error()), log.ERROR)
		return 0, err
	}

	// Query 
	var Id uint64 
	err = stmt.QueryRowContext(ctx, &c.Id).Scan(&Id)
	if err == sql.ErrNoRows {
		log.Log(SOURCE_CONTACT_LIST, fmt.Sprintf("No rows found while reading id %d", Id), log.WARNING)
		return 0, err
	} else if c.Id != Id {
		log.Log(SOURCE_CONTACT_LIST, fmt.Sprintf("Deleted id did not match what was passed. Got: %d, expected: %d.", Id, c.Id), log.ERROR)
		return 0, err
	} else if err != nil {
		log.Log(SOURCE_CONTACT_LIST, fmt.Sprintf("Failed to execute stmt to read contact list for id %d: %s", Id, err.Error()), log.ERROR)
		return 0, err
	}

	return Id, nil
}

// Does not check for Id equality
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

