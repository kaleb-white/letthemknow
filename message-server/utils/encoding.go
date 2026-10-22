package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/kaleb-white/letthemknow/message-server/log"
)

const SOURCE_ENCODER string = "encoder"

// Credit to: https://grafana.com/blog/how-i-write-http-services-in-go-after-13-years/

func Encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		s := fmt.Sprintf("Failed to encode json: %s", err.Error())
		log.Log(SOURCE_ENCODER, s, log.DEBUG)
		return errors.New(s)
	}
	return nil
}

func Decode[T any](r *http.Request) (t T, e error) {
	defer func() {
		if r := recover(); r != nil {
			e = fmt.Errorf("Body decode failed: %s", r)
		}
	}()

	var v T

	if r.Body == nil {
		s := fmt.Sprintf("No body in request for %s", r.URL.RawPath)
		log.Log(SOURCE_ENCODER, s, log.DEBUG)
		return v, errors.New(s)
	}

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		s := fmt.Sprintf("Failed to decode json: %s", err.Error())
		log.Log(SOURCE_ENCODER, s, log.DEBUG)
		return v, errors.New(s)
	}
	return v, nil
}
