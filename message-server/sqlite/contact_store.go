package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/kaleb-white/letthemknow/message-server/log"
	"github.com/kaleb-white/letthemknow/message-server/models"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)

const SQLITE_CONTACT_STORE string = "contact_store"

var CONTACT_PREPARED_SQL = map[string]PreparedSql{
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

type ContactStore struct {
}

var initContactOnce sync.Once

func (_ *ContactStore) Init(ctx context.Context) error {
	var err error
	initContactOnce.Do(func () {
		err = InitializeTable(ctx, CONTACT_TABLEDEF, "contact table")
	})
	if err != nil {
		return err
	}

	return nil
}

func (_ *ContactStore) Read(ctx context.Context, Id uint64) (models.Contact, utils.StoreError) {
	c := models.Contact{}

	// Connect
	conn, err := GetConnection()
	if err != nil {
		e := fmt.Sprintf("Failed to connect to read contact with id %d", Id)
		log.Log(SQLITE_CONTACT_STORE, e, log.INFO)
		return c, utils.NewStoreError(e, http.StatusInternalServerError) 
	}
	defer conn.Close()
	
	// Use statement singleton
	preparedSql := CONTACT_PREPARED_SQL["Read"]
	stmt, err := preparedSql.GetPreparedStatement(ctx, conn)
	if err != nil {
		e := fmt.Sprintf("Failed to prepare stmt to read contact for id %d: %s", Id, err.Error())
		log.Log(SQLITE_CONTACT_STORE, e, log.ERROR)
		return c, utils.NewStoreError(e, http.StatusInternalServerError) 
	}

	// Query 
	var CreatedAtDateTime, LastUpdatedAtDateTime, ListMembership []byte 
	err = stmt.QueryRowContext(ctx, Id).Scan(&c.Id, &c.Phone, &c.Phone2, &c.Phone3, &c.FirstName, &c.LastName, &c.FullName, &ListMembership, &CreatedAtDateTime, &c.CreatedBy, &LastUpdatedAtDateTime, &c.LastUpdatedBy, &c.Org)
	if err == sql.ErrNoRows {
		e := fmt.Sprintf("%d not found", Id)
		log.Log(SQLITE_CONTACT_STORE, e, log.WARNING)
		return c, utils.NewStoreError(e, http.StatusNotFound) 
	} else if err != nil {
		e := fmt.Sprintf("Failed to execute stmt to read contact for id %d: %s", Id, err.Error())
		log.Log(SQLITE_CONTACT_STORE, e, log.ERROR)
		return c, utils.NewStoreError(e, http.StatusInternalServerError) 
	}

	// Map
	err =	c.CreatedAtDateTime.UnmarshalBinary(CreatedAtDateTime)
	if err != nil {
		log.Log(SQLITE_CONTACT_STORE, fmt.Sprintf("Failed to unmarshal CreatedAtDateTime. Raw: %s", CreatedAtDateTime), log.WARNING)
	}
	err = c.LastUpdatedAtDateTime.UnmarshalBinary(LastUpdatedAtDateTime)
	if err != nil {
		log.Log(SQLITE_CONTACT_STORE, fmt.Sprintf("Failed to unmarshal LastUpdatedAtDateTime. Raw: %s", LastUpdatedAtDateTime), log.WARNING)
	}
	err = json.Unmarshal(ListMembership, &c.ListMembership)
	if err != nil {
		log.Log(SQLITE_CONTACT_STORE, fmt.Sprintf("Failed to unmarshal ListMembership. Raw: %s", LastUpdatedAtDateTime), log.WARNING)
	}

	return c, nil
}

func (_ *ContactStore) Write(ctx context.Context, c *models.Contact) (uint64, utils.StoreError) {
	// Connect
	conn, err := GetConnection()
	if err != nil {
		e := fmt.Sprintf("Failed to connect to write contact with id %d", c.Id)
		log.Log(SQLITE_CONTACT_STORE, e, log.INFO)
		return 0, utils.NewStoreError(e, http.StatusInternalServerError) 
	}
	defer conn.Close()
	
	// Use statement singleton: zero value Id means new
	var preparedSql PreparedSql	
	if c.Id == 0 {
		preparedSql = CONTACT_PREPARED_SQL["WriteNew"]
	} else {
		preparedSql = CONTACT_PREPARED_SQL["WriteExisting"]
	}
	stmt, err := preparedSql.GetPreparedStatement(ctx, conn)
	if err != nil {
		e := fmt.Sprintf("Failed to prepare stmt to write contact for id %d: %s", c.Id, err.Error())
		log.Log(SQLITE_CONTACT_STORE, e, log.ERROR)
		return 0, utils.NewStoreError(e, http.StatusInternalServerError) 
	}

	// Map
	CreatedAtDateTime, err := c.CreatedAtDateTime.MarshalBinary()
	if err != nil {
		e := fmt.Sprintf("Failed to marshal CreatedAtDateTime. Value of CreatedAtDateTime: %s", c.CreatedAtDateTime)
		log.Log(SQLITE_CONTACT_STORE, e, log.ERROR)
		return 0, utils.NewStoreError(e, http.StatusInternalServerError) 
	}
	LastUpdatedAtDateTime, err := time.Now().MarshalBinary()
	if err != nil {
		e := fmt.Sprintf("Failed to marshal LastUpdatedAtDateTime. Value of LastUpdatedAtDateTime: %s", c.LastUpdatedAtDateTime)
		log.Log(SQLITE_CONTACT_STORE, e, log.ERROR)
		return 0, utils.NewStoreError(e, http.StatusInternalServerError) 
	}
	ListMembership, err := json.Marshal(c.ListMembership)
	if err != nil {
		e := fmt.Sprintf("Failed to marshal ListMembership. Len of ListMembership: %d", len(c.ListMembership))
		log.Log(SQLITE_CONTACT_STORE, e, log.ERROR)
		return 0, utils.NewStoreError(e, http.StatusInternalServerError) 
	}

	// Query
	var Id uint64
	if c.Id == 0 {
		err = stmt.QueryRowContext(ctx, &c.Phone, &c.Phone2, &c.Phone3, &c.FirstName, &c.LastName, &c.FullName, &ListMembership, &CreatedAtDateTime, &c.CreatedBy, &LastUpdatedAtDateTime, &c.LastUpdatedBy, &c.Org).Scan(&Id)
	} else {
		err = stmt.QueryRowContext(ctx, &c.Id, &c.Phone, &c.Phone2, &c.Phone3, &c.FirstName, &c.LastName, &c.FullName, &ListMembership, &CreatedAtDateTime, &c.CreatedBy, &LastUpdatedAtDateTime, &c.LastUpdatedBy, &c.Org).Scan(&Id)
	}

	if err == sql.ErrNoRows {
		e := fmt.Sprintf("No rows returned after writing contact for phone %d", c.Phone)
		log.Log(SQLITE_CONTACT_STORE, e, log.WARNING)
		return Id, utils.NewStoreError(e, http.StatusNotFound) 
	} else if err != nil {
		e := fmt.Sprintf("Failed to execute stmt to write contact for phone %d: %s", c.Phone, err.Error())
		log.Log(SQLITE_CONTACT_STORE, e, log.ERROR)
		return Id, utils.NewStoreError(e, http.StatusInternalServerError) 
	}

	c.Id = Id

	return Id, nil
}

func (_ *ContactStore) Delete(ctx context.Context, Id uint64) (uint64, utils.StoreError) {
	// Connect
	conn, err := GetConnection()
	if err != nil {
		e := fmt.Sprintf("Failed to connect to delete contact with id %d", Id)
		log.Log(SQLITE_CONTACT_STORE, e, log.INFO)
		return 0, utils.NewStoreError(e, http.StatusInternalServerError) 
	}
	defer conn.Close()
	
	// Use statement singleton
	preparedSql := CONTACT_PREPARED_SQL["Delete"]
	stmt, err := preparedSql.GetPreparedStatement(ctx, conn)
	if err != nil {
		e := fmt.Sprintf("Failed to prepare stmt to delete contact for id %d: %s", Id, err.Error()
		log.Log(SQLITE_CONTACT_STORE, )e, log.ERROR)
		return 0, utils.NewStoreError(e, http.StatusInternalServerError) 
	}

	// Query 
	var reportedId uint64 
	err = stmt.QueryRowContext(ctx, &Id).Scan(&reportedId)
	if err == sql.ErrNoRows {
		e := fmt.Sprintf("%d not found", Id)
		log.Log(SQLITE_CONTACT_STORE, e, log.WARNING)
		return 0, utils.NewStoreError(e, http.StatusNotFound) 
	} else if Id != reportedId {
		e := fmt.Sprintf("Deleted id did not match what was passed. Got: %d, expected: %d.", Id, reportedId)
		log.Log(SQLITE_CONTACT_STORE, e, log.ERROR)
		return 0, utils.NewStoreError(e, http.StatusInternalServerError) 
	} else if err != nil {
		e := fmt.Sprintf("Failed to execute stmt to read contact for id %d: %s", Id, err.Error()
		log.Log(SQLITE_CONTACT_STORE, )e, log.ERROR)
		return 0, utils.NewStoreError(e, http.StatusInternalServerError) 
	}

	return Id, nil
}
