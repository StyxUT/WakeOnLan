package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInvalidJson(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/wol", strings.NewReader(`invalid`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	wolHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestValidMacRequest(t *testing.T) {
	reqBody := `{"mac":"aa:bb:cc:dd:ee:ff"}`
	req := httptest.NewRequest(http.MethodPost, "/wol", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	wolHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestInvalidMacRequest(t *testing.T) {
	reqBody := `{"mac":"cc:dd:ee:ff"}`
	req := httptest.NewRequest(http.MethodPost, "/wol", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	wolHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}
