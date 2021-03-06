package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePairDevice(t *testing.T) {
	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(Pair{DeviceID: 1234, UserID: 9876})

	req := httptest.NewRequest(http.MethodPost, "/pair-device", payload)
	rec := httptest.NewRecorder()

	// #1
	// origin := createPairDevice
	// defer func() {
	// 	createPairDevice = origin
	// }()
	// createPairDevice = func(p Pair) error {
	// 	log.Printf("connected to fake database!\n")
	// 	return nil
	// }

	handler := PairDeviceHandler(CreatePairDeviceFunc(func(p Pair) error {
		return nil
	}))
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected %d, but got %d", http.StatusOK, rec.Code)
	}

	expected := `{"status":"active"}`
	if rec.Body.String() != expected {
		t.Errorf("expected %s, but got %s", expected, rec.Body.String())
	}
}
