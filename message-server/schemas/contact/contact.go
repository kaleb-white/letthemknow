package contact

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/kaleb-white/letthemknow/message-server/log"
	"github.com/kaleb-white/letthemknow/message-server/utils"
	"github.com/kaleb-white/letthemknow/message-server/db"
)

const SOURCE_CONTACT string = "contact"

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

var CONTACT_PREPARED_SQL = map[string]db.PreparedSql{
	"Read": {
		OperationName: "contact_read",
		RawSql: CONTACT_READ,
	},
	"WriteNew": {
		OperationName: "contact_write",
		RawSql: CONTACT_WRITE_NEW,
	},
	"WriteExisting": {
		OperationName: "contact_write",
		RawSql: CONTACT_WRITE_EXISTING,
	},
	"Delete": {
		OperationName: "contact_delete",
		RawSql: CONTACT_DELETE,
	},
}

func (c *Contact) Read(ctx context.Context, Id uint64) (*Contact, error) {
	// Connect
	conn, err := db.GetConnection()
	if err != nil {
		log.Log(SOURCE_CONTACT, fmt.Sprintf("Failed to connect to read contact with id %d", Id), log.INFO)
		return c, err
	}
	defer conn.Close()
	
	// Use statement singleton
	preparedSql := CONTACT_PREPARED_SQL["Read"]
	stmt, err := preparedSql.GetPreparedStatement(ctx, conn)
	if err != nil {
		log.Log(SOURCE_CONTACT, fmt.Sprintf("Failed to prepare stmt to read contact for id %d: %s", Id, err.Error()), log.ERROR)
		return c, err
	}

	// Query 
	var CreatedAtDateTime, LastUpdatedAtDateTime []byte 
	err = stmt.QueryRowContext(ctx, Id).Scan(&c.Id, &c.Phone, &c.Phone2, &c.Phone3, &c.FirstName, &c.LastName, &c.FullName, &CreatedAtDateTime, &c.CreatedBy, &LastUpdatedAtDateTime, &c.LastUpdatedBy, &c.Org)
	if err == sql.ErrNoRows {
		log.Log(SOURCE_CONTACT, fmt.Sprintf("No rows found while reading id %d", Id), log.WARNING)
		return c, err
	} else if err != nil {
		log.Log(SOURCE_CONTACT, fmt.Sprintf("Failed to execute stmt to read contact for id %d: %s", Id, err.Error()), log.ERROR)
		return c, err
	}

	// Map
	err =	c.CreatedAtDateTime.UnmarshalBinary(CreatedAtDateTime)
	if err != nil {
		log.Log(SOURCE_CONTACT, fmt.Sprintf("Failed to unmarshal CreatedAtDateTime. Raw: %s", CreatedAtDateTime), log.WARNING)
	}
	err = c.LastUpdatedAtDateTime.UnmarshalBinary(LastUpdatedAtDateTime)
	if err != nil {
		log.Log(SOURCE_CONTACT, fmt.Sprintf("Failed to unmarshal LastUpdatedAtDateTime. Raw: %s", LastUpdatedAtDateTime), log.WARNING)
	}

	return c, nil
}

func (c *Contact) Write(ctx context.Context) (uint64, error) {
	// Connect
	conn, err := db.GetConnection()
	if err != nil {
		log.Log(SOURCE_CONTACT, fmt.Sprintf("Failed to connect to write contact with id %d", c.Id), log.INFO)
		return 0, err
	}
	defer conn.Close()
	
	// Use statement singleton: zero value Id means new
	var preparedSql db.PreparedSql	
	if c.Id == 0 {
		preparedSql = CONTACT_PREPARED_SQL["WriteNew"]
	} else {
		preparedSql = CONTACT_PREPARED_SQL["WriteExisting"]
	}
	stmt, err := preparedSql.GetPreparedStatement(ctx, conn)
	if err != nil {
		log.Log(SOURCE_CONTACT, fmt.Sprintf("Failed to prepare stmt to write contact for id %d: %s", c.Id, err.Error()), log.ERROR)
		return 0, err
	}

	// Map
	CreatedAtDateTime, err := c.CreatedAtDateTime.MarshalBinary()
	if err != nil {
		log.Log(SOURCE_CONTACT, fmt.Sprintf("Failed to marshal CreatedAtDateTime. Value of CreatedAtDateTime: %s", c.CreatedAtDateTime), log.ERROR)
		return 0, err
	}
	LastUpdatedAtDateTime, err := time.Now().MarshalBinary()
	if err != nil {
		log.Log(SOURCE_CONTACT, fmt.Sprintf("Failed to marshal LastUpdatedAtDateTime. Value of LastUpdatedAtDateTime: %s", c.LastUpdatedAtDateTime), log.ERROR)
		return 0, err
	}

	// Query
	var Id uint64
	if c.Id == 0 {
		err = stmt.QueryRowContext(ctx, &c.Phone, &c.Phone2, &c.Phone3, &c.FirstName, &c.LastName, &c.FullName, &CreatedAtDateTime, &c.CreatedBy, &LastUpdatedAtDateTime, &c.LastUpdatedBy, &c.Org).Scan(&Id)
	} else {
		err = stmt.QueryRowContext(ctx, &c.Id, &c.Phone, &c.Phone2, &c.Phone3, &c.FirstName, &c.LastName, &c.FullName, &CreatedAtDateTime, &c.CreatedBy, &LastUpdatedAtDateTime, &c.LastUpdatedBy, &c.Org).Scan(&Id)
	}

	if err == sql.ErrNoRows {
		log.Log(SOURCE_CONTACT, fmt.Sprintf("No rows returned after writing contact for phone %d", c.Phone), log.WARNING)
		return Id, err
	} else if err != nil {
		log.Log(SOURCE_CONTACT, fmt.Sprintf("Failed to execute stmt to write contact for phone %d: %s", c.Phone, err.Error()), log.ERROR)
		return Id, err
	}

	c.Id = Id

	return Id, nil
}

func (c *Contact) Delete(ctx context.Context) (uint64, error) {
	// Connect
	conn, err := db.GetConnection()
	if err != nil {
		log.Log(SOURCE_CONTACT, fmt.Sprintf("Failed to connect to delete contact with id %d", c.Id), log.INFO)
		return 0, err
	}
	defer conn.Close()
	
	// Use statement singleton
	preparedSql := CONTACT_PREPARED_SQL["Delete"]
	stmt, err := preparedSql.GetPreparedStatement(ctx, conn)
	if err != nil {
		log.Log(SOURCE_CONTACT, fmt.Sprintf("Failed to prepare stmt to delete contact for id %d: %s", c.Id, err.Error()), log.ERROR)
		return 0, err
	}

	// Query 
	var Id uint64 
	err = stmt.QueryRowContext(ctx, &c.Id).Scan(&Id)
	if err == sql.ErrNoRows {
		log.Log(SOURCE_CONTACT, fmt.Sprintf("No rows found while reading id %d", Id), log.WARNING)
		return 0, err
	} else if c.Id != Id {
		log.Log(SOURCE_CONTACT, fmt.Sprintf("Deleted id did not match what was passed. Got: %d, expected: %d.", Id, c.Id), log.ERROR)
		return 0, err
	} else if err != nil {
		log.Log(SOURCE_CONTACT, fmt.Sprintf("Failed to execute stmt to read contact for id %d: %s", Id, err.Error()), log.ERROR)
		return 0, err
	}

	return Id, nil
}

// Does not check for Id equality
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

