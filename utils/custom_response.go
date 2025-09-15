package utils

import (
	"encoding/json"
	"net/http"
)

type CustomResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Data       any    `json:"data,omitempty"`
}

// func CustomError(w http.ResponseWriter, statusCode int, message string) {
// 	resp := CustomResponse{
// 		Message:    message,
// 		StatusCode: statusCode,
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(statusCode)
// 	if err := json.NewEncoder(w).Encode(resp); err != nil {
// 		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
// 	}
// }
// func SendCustomResponse(w http.ResponseWriter, statusCode int, message string, data any) {
// 	resp := CustomResponse{
// 		Message:    message,
// 		StatusCode: statusCode,
// 		Data:       data,
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(statusCode)
// 	if err := json.NewEncoder(w).Encode(resp); err != nil {
// 		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
// 	}
// }

func CustomResponseSender(w http.ResponseWriter, statusCode int, message string, data ...any) {
	var resp CustomResponse
	if len(data) != 0 {
		resp = CustomResponse{
			Message:    message,
			StatusCode: statusCode,
			Data:       data[0],
		}
	} else {
		resp = CustomResponse{
			Message:    message,
			StatusCode: statusCode,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

}
