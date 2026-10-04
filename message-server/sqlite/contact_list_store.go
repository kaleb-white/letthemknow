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

const SQLITE_CONTACT_LIST string = "contact_list"

var CONTACT_LIST_PREPARED_SQL = map[string]PreparedSql{
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

type ContactListStore struct {
}

var onceContactList sync.Once

func (_ *ContactListStore) Init(ctx context.Context) error {
	var err error
	onceContactList.Do(func () {
		err = InitializeTable(ctx, CONTACT_LIST_TABLEDEF, "contact list table")
	})
	if err != nil {
		return err
	}

	return nil
}

func (_ *ContactListStore) Read(ctx context.Context, Id uint64) (models.ContactList, utils.StoreError) {
	c := models.ContactList{} 

	// Connect
	conn, err := GetConnection()
	if err != nil {
		e := fmt.Sprintf("Failed to connect to read contact list with id %d", Id)
		log.Log(SQLITE_CONTACT_LIST, e, log.INFO)
		return c, utils.NewStoreError(e, http.StatusInternalServerError)
	}
	defer conn.Close()
	
	// Use statement singleton
	preparedSql := CONTACT_LIST_PREPARED_SQL["Read"]
	stmt, err := preparedSql.GetPreparedStatement(ctx, conn)
	if err != nil {
		e := fmt.Sprintf("Failed to prepare stmt to read contact list for id %d: %s", Id, err.Error())
		log.Log(SQLITE_CONTACT_LIST, e, log.ERROR)
		return c, utils.NewStoreError(e, http.StatusInternalServerError)
	}

	// Query 
	var CreatedAtDateTime, LastUpdatedAtDateTime, ContactList []byte 
	err = stmt.QueryRowContext(ctx, Id).Scan(&c.Id, &c.Name, &c.Description, &ContactList, &CreatedAtDateTime, &c.CreatedBy, &LastUpdatedAtDateTime, &c.LastUpdatedBy)
	if err == sql.ErrNoRows {
		e := fmt.Sprintf("%d not found", Id)
		log.Log(SQLITE_CONTACT_LIST, e, log.WARNING)
		return c, utils.NewStoreError(e, http.StatusNotFound)
	} else if err != nil {
		e := fmt.Sprintf("Failed to execute stmt to read contact list for id %d: %s", Id, err.Error())
		log.Log(SQLITE_CONTACT_LIST, e, log.ERROR)
		return c, utils.NewStoreError(e, http.StatusInternalServerError)
	}

	// Map
	err =	c.CreatedAtDateTime.UnmarshalBinary(CreatedAtDateTime)
	if err != nil {
		log.Log(SQLITE_CONTACT_LIST, fmt.Sprintf("Failed to unmarshal CreatedAtDateTime. Raw: %s", CreatedAtDateTime), log.WARNING)
	}
	err = c.LastUpdatedAtDateTime.UnmarshalBinary(LastUpdatedAtDateTime)
	if err != nil {
		log.Log(SQLITE_CONTACT_LIST, fmt.Sprintf("Failed to unmarshal LastUpdatedAtDateTime. Raw: %s", LastUpdatedAtDateTime), log.WARNING)
	}
	err = json.Unmarshal(ContactList, &c.ContactList)
	if err != nil {
		log.Log(SQLITE_CONTACT_LIST, fmt.Sprintf("Failed to unmarshal ContactList. Raw: %s. Err: %s", LastUpdatedAtDateTime, err.Error()), log.WARNING)
	}

	return c, nil
}

func (_ *ContactListStore) Write(ctx context.Context, c *models.ContactList) (uint64, utils.StoreError) {
	// Connect
	conn, err := GetConnection()
	if err != nil {
		e := fmt.Sprintf("Failed to connect to write contact list with id %d", c.Id)
		log.Log(SQLITE_CONTACT_LIST, e, log.INFO)
		return 0, utils.NewStoreError(e, http.StatusInternalServerError) 
	}
	defer conn.Close()
	
	// Use statement singleton: zero value Id means new
	var preparedSql PreparedSql	
	if c.Id == 0 {
		preparedSql = CONTACT_LIST_PREPARED_SQL["WriteNew"]
	} else {
		preparedSql = CONTACT_LIST_PREPARED_SQL["WriteExisting"]
	}
	stmt, err := preparedSql.GetPreparedStatement(ctx, conn)
	if err != nil {
		e := fmt.Sprintf("Failed to prepare stmt to write contact list for id %d: %s", c.Id, err.Error())
		log.Log(SQLITE_CONTACT_LIST, e, log.ERROR)
		return 0, utils.NewStoreError(e, http.StatusInternalServerError) 
	}

	// Map
	CreatedAtDateTime, err := c.CreatedAtDateTime.MarshalBinary()
	if err != nil {
		e := fmt.Sprintf("Failed to marshal CreatedAtDateTime. Value of CreatedAtDateTime: %s", c.CreatedAtDateTime)
		log.Log(SQLITE_CONTACT_LIST, e, log.ERROR)
		return 0, utils.NewStoreError(e, http.StatusInternalServerError) 
	}
	LastUpdatedAtDateTime, err := time.Now().MarshalBinary()
	if err != nil {
		e := fmt.Sprintf("Failed to marshal LastUpdatedAtDateTime. Value of LastUpdatedAtDateTime: %s", c.LastUpdatedAtDateTime)
		log.Log(SQLITE_CONTACT_LIST, e, log.ERROR)
		return 0, utils.NewStoreError(e, http.StatusInternalServerError) 
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
		e := fmt.Sprintf("No rows returned after writing contact list %s", c.Name)
		log.Log(SQLITE_CONTACT_LIST, e, log.WARNING)
		return Id, utils.NewStoreError(e, http.StatusInternalServerError) 
	} else if err != nil {
		e := fmt.Sprintf("Failed to execute stmt to write for contact list %s: %s", c.Name, err.Error())
		log.Log(SQLITE_CONTACT_LIST, e, log.ERROR)
		return Id, utils.NewStoreError(e, http.StatusInternalServerError) 
	}

	c.Id = Id

	return Id, nil
}

func (_ *ContactListStore) Delete(ctx context.Context, Id uint64) (uint64, utils.StoreError) {
	// Connect
	conn, err := GetConnection()
	if err != nil {
		e := fmt.Sprintf("Failed to connect to delete ContactList with id %d", Id)
		log.Log(SQLITE_CONTACT_LIST, e, log.INFO)
		return 0, utils.NewStoreError(e, http.StatusInternalServerError) 
	}
	defer conn.Close()
	
	// Use statement singleton
	preparedSql := CONTACT_LIST_PREPARED_SQL["Delete"]
	stmt, err := preparedSql.GetPreparedStatement(ctx, conn)
	if err != nil {
		e := fmt.Sprintf("Failed to prepare stmt to delete contact list for id %d: %s", Id, err.Error())
		log.Log(SQLITE_CONTACT_LIST, e, log.ERROR)
		return 0, utils.NewStoreError(e, http.StatusInternalServerError) 
	}

	// Query 
	var reportedId uint64 
	err = stmt.QueryRowContext(ctx, &Id).Scan(&reportedId)
	if err == sql.ErrNoRows {
		e := fmt.Sprintf("%d not found", Id)
		log.Log(SQLITE_CONTACT_LIST, e, log.WARNING)
		return 0, utils.NewStoreError(e, http.StatusNotFound) 
	} else if Id != reportedId {
		e := fmt.Sprintf("Deleted id did not match what was passed. Got: %d, expected: %d.", Id, reportedId)
		log.Log(SQLITE_CONTACT_LIST, e, log.ERROR)
		return 0, utils.NewStoreError(e, http.StatusInternalServerError) 
	} else if err != nil {
		e := fmt.Sprintf("Failed to execute stmt to read contact list for id %d: %s", Id, err.Error())
		log.Log(SQLITE_CONTACT_LIST, e, log.ERROR)
		return 0, utils.NewStoreError(e, http.StatusInternalServerError) 
	}

	return Id, nil
}
