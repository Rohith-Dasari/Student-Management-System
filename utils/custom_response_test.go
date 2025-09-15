package utils_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sms/utils"
	"testing"
)

// func TestCustomError(t *testing.T) {
// 	rr := httptest.NewRecorder()
// 	utils.CustomError(rr, http.StatusBadRequest, "test error")

// 	if rr.Code != http.StatusBadRequest {
// 		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
// 	}

// 	var resp utils.CustomResponse
// 	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
// 		t.Fatalf("failed to decode response: %v", err)
// 	}

// 	if resp.Message != "test error" {
// 		t.Errorf("expected message %q, got %q", "test error", resp.Message)
// 	}

// 	if resp.StatusCode != http.StatusBadRequest {
// 		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
// 	}
// }

// func TestSendCustomResponse(t *testing.T) {
// 	rr := httptest.NewRecorder()

// 	data := map[string]string{"key": "value"}
// 	utils.SendCustomResponse(rr, http.StatusOK, "success", data)

// 	if rr.Code != http.StatusOK {
// 		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
// 	}

// 	var resp utils.CustomResponse
// 	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
// 		t.Fatalf("failed to decode response: %v", err)
// 	}

// 	if resp.Message != "success" {
// 		t.Errorf("expected message %q, got %q", "success", resp.Message)
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
// 	}

// 	if val, ok := resp.Data.(map[string]any); ok {
// 		if val["key"] != "value" {
// 			t.Errorf("expected data key 'value', got %v", val["key"])
// 		}
// 	} else {
// 		t.Errorf("expected data to be map[string]interface{}, got %T", resp.Data)
// 	}
// }

func TestCustomResponseSender(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		message      string
		data         []any
		expectedBody utils.CustomResponse
	}{
		{
			name:       "Response with data",
			statusCode: http.StatusOK,
			message:    "ok",
			data:       []any{"some data"},
			expectedBody: utils.CustomResponse{
				Message:    "ok",
				StatusCode: http.StatusOK,
				Data:       "some data",
			},
		},
		{
			name:       "Response without data",
			statusCode: http.StatusBadRequest,
			message:    "bad request",
			data:       nil,
			expectedBody: utils.CustomResponse{
				Message:    "bad request",
				StatusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			utils.CustomResponseSender(rr, tt.statusCode, tt.message, tt.data...)

			if rr.Code != tt.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.statusCode)
			}
			if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("handler returned wrong Content-Type: got %v want %v", contentType, "application/json")
			}
			var actualBody utils.CustomResponse
			if err := json.NewDecoder(rr.Body).Decode(&actualBody); err != nil {
				t.Fatalf("could not decode response body: %v", err)
			}
			if !reflect.DeepEqual(actualBody, tt.expectedBody) {
				t.Errorf("handler returned unexpected body:\nGot:  %v\nWant: %v", actualBody, tt.expectedBody)
			}
		})
	}
}
