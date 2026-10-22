package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func Getenv(string) string {
	return "1000"
}

func GenerateRequest(t *testing.T, method string, url string, body any) *http.Request {
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(body)
	if err != nil {
		t.Fatal(err.Error())
	}

	rdr := bytes.NewReader(b.Bytes())
	req, err := http.NewRequest(method, url, rdr)
	if err != nil {
		t.Fatal(err.Error())
	}
	return req
}
