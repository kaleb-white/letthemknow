package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/kaleb-white/letthemknow/message-server/log"
	"github.com/kaleb-white/letthemknow/message-server/models"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)

const SQLITE_CONTACT_DETAIL_STORE = "contact_detail_STORE"

var CONTACT_DETAIL_PREPARED_SQL = map[string]PreparedSql{
	"Read": {
		OperationName: "contact_details_read",
		RawSql: CONTACT_DETAIL_READ,
	},
	"Delete": {
		OperationName: "contact_details_delete",
		RawSql: CONTACT_DETAIL_DELETE,
	},
}

type ContactDetailStore struct {}

func (_ *ContactDetailStore) Init(ctx context.Context) error {
	contactStore := ContactStore{}
	return contactStore.Init(ctx)
}

func (_ *ContactDetailStore) Read(ctx context.Context, Id uint64) (models.ContactDetail, utils.StoreError) {
	c := models.ContactDetail{}

	// Connect
	conn, err := GetConnection()
	if err != nil {
		e := fmt.Sprintf("Failed to connect to read contact details with id %d: %s", Id, err.Error())
		log.Log(SQLITE_CONTACT_DETAIL_STORE, e, log.ERROR)
		return c, utils.NewStoreError(e, http.StatusInternalServerError)
	}
	defer conn.Close()

	// Use statement singleton
	preparedSql := CONTACT_DETAIL_PREPARED_SQL["Read"]
	stmt, err := preparedSql.GetPreparedStatement(ctx, conn)
	if err != nil {
		e := fmt.Sprintf("Failed to prepare stmt to read contact details for id %d: %s", Id, err.Error())
		log.Log(SQLITE_CONTACT_DETAIL_STORE, e, log.ERROR)
		return c, utils.NewStoreError(e, http.StatusInternalServerError)
	}

	// Query 
	var LastUpdatedAtDateTime []byte 
	err = stmt.QueryRowContext(ctx, Id).Scan(&c.Id, &c.Phone, &c.FirstName, &c.LastName, &LastUpdatedAtDateTime, &c.LastUpdatedBy, &c.Org)

	if err == sql.ErrNoRows {
		e := fmt.Sprintf("%d not found", Id)
		log.Log(SQLITE_CONTACT_DETAIL_STORE, e, log.DEBUG)
		return c, utils.NewStoreError(e, http.StatusNotFound)
	} else if err != nil {
		e := fmt.Sprintf("Failed to execute stmt to read contact details for id %d: %s", Id, err.Error())
		log.Log(SQLITE_CONTACT_DETAIL_STORE, e, log.ERROR)
		return c, utils.NewStoreError(e, http.StatusInternalServerError) 
	}

	// Map
	err = c.LastUpdatedAtDateTime.UnmarshalBinary(LastUpdatedAtDateTime)
	if err != nil {
		log.Log(SQLITE_CONTACT_DETAIL_STORE, fmt.Sprintf("Failed to unmarshal LastUpdatedAtDateTime. Raw: %s", LastUpdatedAtDateTime), log.WARNING)
	}

	return c, nil
}

func (_ *ContactDetailStore) Delete(ctx context.Context, Id uint64) (uint64, utils.StoreError) {
	// Connect
	conn, err := GetConnection()
	if err != nil {
		e := fmt.Sprintf("Failed to connect to delete contact details with id %d", Id)
		log.Log(SQLITE_CONTACT_DETAIL_STORE, e, log.INFO)
		return 0, utils.NewStoreError(e, http.StatusInternalServerError)
	}
	defer conn.Close()

	// Use statement singleton
	preparedSql := CONTACT_DETAIL_PREPARED_SQL["Delete"]
	stmt, err := preparedSql.GetPreparedStatement(ctx, conn)
	if err != nil {
		e := fmt.Sprintf("Failed to prepare stmt to delete contact details for id %d: %s", Id, err.Error())
		log.Log(SQLITE_CONTACT_DETAIL_STORE, e, log.ERROR)
		return 0, utils.NewStoreError(e, http.StatusInternalServerError)
	}

	// Query 
	var reportedId uint64 
	err = stmt.QueryRowContext(ctx, &Id).Scan(&reportedId)
	if err == sql.ErrNoRows {
		e := fmt.Sprintf("%d not found", Id)
		log.Log(SQLITE_CONTACT_DETAIL_STORE, e, log.WARNING)
		return 0, utils.NewStoreError(e, http.StatusNotFound)
	} else if Id != reportedId {
		e := fmt.Sprintf("Deleted id did not match what was passed. Got: %d, expected: %d.", Id, reportedId)
		log.Log(SQLITE_CONTACT_DETAIL_STORE, e, log.ERROR)
		return 0, utils.NewStoreError(e, http.StatusInternalServerError)
	} else if err != nil {
		e := fmt.Sprintf("Failed to execute stmt to read contact for id %d: %s", Id, err.Error())
		log.Log(SQLITE_CONTACT_DETAIL_STORE, e, log.ERROR)
		return 0, utils.NewStoreError(e, http.StatusInternalServerError)
	}

	return Id, nil
}
