package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kaleb-white/letthemknow/message-server/log"
	"github.com/kaleb-white/letthemknow/message-server/models"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)

func GenericReadHandler[
	N any, 
	S models.ReadEnabledStore[N],
](
	store S, 
	getenv func(string) string, 
	logSource string, 
	typeName string,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// This should already be checked by HasSlugs middleware
			id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)

			// Perform request to store
			obj, err := ReadFromStoreWithContext(id, store, getenv)
			
			// Parse result
			if err != nil {
				e := fmt.Sprintf("Failure while reading %s: %s", typeName, err.Error())
				log.Log(logSource, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(err.Error(), err.Status(), w, r)
				return
			}

			encodeErr := utils.Encode(w, r, http.StatusOK, obj)
			if encodeErr != nil {
				e := fmt.Sprintf("Failed to encode %s for read: %s", typeName, encodeErr.Error())
				log.Log(logSource, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(e, http.StatusInternalServerError, w, r)
			}

			log.Log(logSource, utils.AttributeLog(fmt.Sprintf("Read %s with id %d", typeName, id), r), log.INFO)
		},
	)
}

func GenericDeleteHandler[
	N any, 
	S models.DeleteEnabledStore[N],
](
	store S, 
	getenv func(string) string, 
	logSource string, 
	typeName string,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// This should already be checked by HasSlugs middleware
			id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)

			// Perform request to store
			timeoutDur := utils.ConfiguredDuration("SQLITE_DELETE_TIMEOUT", getenv)
			ctx, cancel := context.WithTimeout(context.Background(), timeoutDur)
			defer cancel()
			obj, err := store.Delete(ctx, id)

			// Parse result
			if err != nil {
				e := fmt.Sprintf("Failure while reading %s: %s", typeName, err.Error())
				log.Log(logSource, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(err.Error(), err.Status(), w, r)
				return
			}

			encodeErr := utils.Encode(w, r, http.StatusOK, obj)
			if encodeErr != nil {
				e := fmt.Sprintf("Failed to encode %s during delete", typeName)
				log.Log(logSource, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(e, http.StatusInternalServerError, w, r)
			}

			log.Log(logSource, utils.AttributeLog(fmt.Sprintf("Deleted %s with id %d", typeName, id), r), log.INFO)
		},
	)

}

func GenericWriteHandler[
	N any,
	M models.IdentifiableValidatable[N],
	S models.WriteEnabledStore[N],
](
	store S,
	getenv func(string) string,
	logSource string,
	typeName string,
	preValidation func(M, *http.Request),
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Decode JSON body
			obj, decodeErr := utils.Decode[N](r)
			if decodeErr != nil {
				e := fmt.Sprintf("Failed to decode request body while writing %s: %s", typeName, decodeErr.Error())
				log.Log(logSource, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(e, http.StatusBadRequest, w, r)
				return
			}

			// Run prevalidation
			preValidation(&obj, r)

			// Validate input
			var objPtr M = &obj
			if err, errs := objPtr.Validate(); err {
				e := utils.PrettyPrintErrors(errs)
				log.Log(logSource, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(e, http.StatusBadRequest, w, r)
				return
			}

			Id, err := WriteToStoreWithContext(objPtr, store, getenv)
			if err != nil {
				e := err.Error()
				log.Log(logSource, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(e, err.Status(), w, r)
				return
			}

			resp := models.IdResponse{
				Id: Id,
			}		
			encodeErr := utils.Encode(w, r, http.StatusOK, resp)
			if encodeErr != nil {
				e := "Failed to encode return object"
				log.Log(logSource, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(e, http.StatusInternalServerError, w, r)
			}

			log.Log(logSource, utils.AttributeLog(fmt.Sprintf("Wrote %s with id %d", typeName, Id), r), log.INFO)
		},
	)
}

func GenericWriteByIdHandler[
	N any,
	M models.IdentifiableValidatable[N],
	S models.WriteEnabledStore[N],
](
	store S,
	getenv func(string) string,
	logSource string,
	typeName string,
	preValidation func(*N, *http.Request),
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// This should already be checked by HasSlugs middleware
			id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)

			// Decode JSON body
			obj, decodeErr := utils.Decode[N](r)
			if decodeErr != nil {
				e := fmt.Sprintf("Failed to decode request body while writing %s: %s", typeName, decodeErr.Error())
				log.Log(logSource, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(e, http.StatusBadRequest, w, r)
				return
			}

			// Run prevalidation
			preValidation(&obj, r)

			// Validate input
			var objPtr M = &obj
			if err, errs := objPtr.Validate(); err {
				e := utils.PrettyPrintErrors(errs)
				log.Log(logSource, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(e, http.StatusBadRequest, w, r)
				return
			}

			// Overwrite id for store
			objPtr.SetId(id)

			Id, err := WriteToStoreWithContext(objPtr, store, getenv)
			if err != nil {
				e := err.Error()
				log.Log(logSource, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(e, err.Status(), w, r)
				return
			}

			rObj := models.IdResponse{
				Id: Id,
			}
			encodeErr := utils.Encode(w, r, http.StatusOK, rObj)
			if encodeErr != nil {
				e := "Failed to encode return object"
				log.Log(logSource, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(e, http.StatusInternalServerError, w, r)
			}

			log.Log(logSource, utils.AttributeLog(fmt.Sprintf("Wrote %s with id %d", typeName, Id), r), log.INFO)
		},
	)
}
