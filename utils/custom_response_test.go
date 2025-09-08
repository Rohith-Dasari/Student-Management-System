package utils_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sms/utils"
	"testing"
)

func TestCustomError(t *testing.T) {
	rr := httptest.NewRecorder()
	utils.CustomError(rr, http.StatusBadRequest, "test error")

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	var resp utils.CustomResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Message != "test error" {
		t.Errorf("expected message %q, got %q", "test error", resp.Message)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestSendCustomResponse(t *testing.T) {
	rr := httptest.NewRecorder()

	data := map[string]string{"key": "value"}
	utils.SendCustomResponse(rr, http.StatusOK, "success", data)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp utils.CustomResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Message != "success" {
		t.Errorf("expected message %q, got %q", "success", resp.Message)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if val, ok := resp.Data.(map[string]any); ok {
		if val["key"] != "value" {
			t.Errorf("expected data key 'value', got %v", val["key"])
		}
	} else {
		t.Errorf("expected data to be map[string]interface{}, got %T", resp.Data)
	}
}
