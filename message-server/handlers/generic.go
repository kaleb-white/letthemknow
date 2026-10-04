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

func genericReadHandler[M any, S models.ReadEnabledStore[M]](store *S, getenv func(string) string, logSource string, typeName string) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			s := *store

			// This should already be checked by HasSlugs middleware
			id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)

			// Perform request to store
			timeoutDur := utils.ConfiguredDuration("SQLITE_READ_TIMEOUT", getenv)
			ctx, cancel := context.WithTimeout(context.Background(), timeoutDur)
			defer cancel()
			obj, err := s.Read(ctx, id)

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

func genericDeleteHandler[M any, S models.DeleteEnabledStore[M]](store *S, getenv func(string) string, logSource string, typeName string) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cs := *store

			// This should already be checked by HasSlugs middleware
			id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)

			// Perform request to store
			timeoutDur := utils.ConfiguredDuration("SQLITE_DELETE_TIMEOUT", getenv)
			ctx, cancel := context.WithTimeout(context.Background(), timeoutDur) 
			defer cancel()
			obj, err := cs.Delete(ctx, id)

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

func genericWriteHandler[
	N any, 
	M models.IdentifiableValidatable[N], 
	S models.WriteEnabledStore[N],
](
	store *S, 
	getenv func(string) string, 
	logSource string, 
	typeName string, 
	preValidation func(M, *http.Request),
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			s := *store

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
			if pass, errs := objPtr.Validate(); !pass {
				e := utils.PrettyPrintErrors(errs)
				log.Log(logSource, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(e, http.StatusBadRequest, w, r)	
				return
			}

			timeoutDur := utils.ConfiguredDuration("SQLITE_WRITE_TIMEOUT", getenv)
			ctx, cancel := context.WithTimeout(context.Background(), timeoutDur) 
			defer cancel()
			Id, err := s.Write(ctx, &obj)
			if err != nil {
				e := err.Error()
				log.Log(logSource, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(e, err.Status(), w, r)
				return
			}

			rObj := utils.WriteId{
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

func genericWriteByIdHandler[
	N any, 
	M models.IdentifiableValidatable[N],  
	S models.WriteEnabledStore[N],
](
	store *S, 
	getenv func(string) string, 
	logSource string, 
	typeName string, 
	preValidation func(*N, *http.Request),
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			s := *store

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
			if pass, errs := objPtr.Validate(); !pass {
				e := utils.PrettyPrintErrors(errs)
				log.Log(logSource, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(e, http.StatusBadRequest, w, r)	
				return
			}

			// Overwrite id for store
			objPtr.SetId(id)
			
			timeoutDur := utils.ConfiguredDuration("SQLITE_WRITE_TINEOUT", getenv)
			ctx, cancel := context.WithTimeout(context.Background(), timeoutDur) 
			defer cancel()
			Id, err := s.Write(ctx, objPtr)
			if err != nil {
				e := err.Error()
				log.Log(logSource, utils.AttributeLog(e, r), log.INFO)
				utils.EncodeHttpError(e, err.Status(), w, r)
				return
			}

			rObj := utils.WriteId{
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
