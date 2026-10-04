package utils

import (
	"fmt"
	"net/http"

	"github.com/kaleb-white/letthemknow/message-server/log"
)

type HttpError struct {
	error string
}

const SOURCE_HTTP_ERROR string = "http_error"

func CreateHttpError(error string, status int) HttpError {
	log.Log(SOURCE_HTTP_ERROR, fmt.Sprintf("%d: %s", status, error), log.DEBUG)
	e := HttpError{
		error: error,
	}
	return e
}

func EncodeHttpError(error string, status int, w http.ResponseWriter, r *http.Request) error {
	httpErr := CreateHttpError(error, status)
	if err := Encode(w, r, status, httpErr); err != nil {
		log.Log(SOURCE_HTTP_ERROR, fmt.Sprintf("Failed to encode err: %s", err.Error()), log.DEBUG)
		w.WriteHeader(http.StatusInternalServerError)		
		return err
	}
	return nil
}
