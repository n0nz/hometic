package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePairDevice(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/pair-device", nil)
	rec := httptest.NewRecorder()

	PairDeviceHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected %d, but got %d", http.StatusOK, rec.Code)
	}

	expected := `{"status":"active"}`
	if rec.Body.String() != expected {
		t.Errorf("expected %s, but got %s", expected, rec.Body.String())
	}
}
