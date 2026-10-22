package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kaleb-white/letthemknow/message-server/middleware"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)

var middlewareTestMux http.ServeMux = *http.NewServeMux()

var doesNothingHandler http.Handler = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request){},
)

func TestHasSlugsGoldenPath(t *testing.T) {
	called := false
	recordsCallHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request){
			called = true
		},
	)

	h := middleware.HasSlugs(recordsCallHandler, []string{"id"})

	path := "/TestHasSlugsGoldenPath/1"
	r := utils.GenerateRequest(t, "GET", path, nil)
	rr := httptest.NewRecorder()
	
	middlewareTestMux.Handle("/TestHasSlugsGoldenPath/{id}", h)

	middlewareTestMux.ServeHTTP(rr, r)

	if !called {
		t.Errorf("Middleware did not call handler")
	}
}

func TestHasSlugsNaN(t *testing.T) {
	h := middleware.HasSlugs(doesNothingHandler, []string{"id"})
	
	path := "/TestHasSlugsNaN/abc"
	r := utils.GenerateRequest(t, "GET", path, nil)
	rr := httptest.NewRecorder()

	middlewareTestMux.Handle("/TestHasSlugsNaN/{id}", h)

	middlewareTestMux.ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler did not succeed, got %v wanted %v", status, http.StatusBadRequest)
		return
	}
}

func TestHasSlugsSlugMissing(t *testing.T) {
	h := middleware.HasSlugs(doesNothingHandler, []string{"id"})
	
	path := "/TestHasSlugsSlugMissing/1"
	r := utils.GenerateRequest(t, "GET", path, nil)
	rr := httptest.NewRecorder()

	middlewareTestMux.Handle("/TestHasSlugsSlugMissing/{di}", h)

	middlewareTestMux.ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler did not succeed, got %v wanted %v", status, http.StatusBadRequest)
		return
	}
}
