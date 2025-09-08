package handlers_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"sms/handlers"
// 	"sms/mocks"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// )

// func TestGradeHandler_AddGrade(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	tests := []struct {
// 		name           string
// 		method         string
// 		body           any
// 		mockSetup      func(mock *mocks.MockGradeServiceI)
// 		expectedStatus int
// 	}{
// 		{
// 			name:   "success",
// 			method: http.MethodPost,
// 			body: map[string]interface{}{
// 				"studentID": "1",
// 				"subjectID": "sub1",
// 				"semester":  1,
// 				"grade":     95,
// 			},
// 			mockSetup: func(mock *mocks.MockGradeServiceI) {
// 				mock.EXPECT().AddGrades("1", "sub1", 95, 1).Return(nil)
// 			},
// 			expectedStatus: http.StatusCreated,
// 		},
// 		{
// 			name:           "invalid method",
// 			method:         http.MethodGet,
// 			body:           nil,
// 			mockSetup:      func(mock *mocks.MockGradeServiceI) {},
// 			expectedStatus: http.StatusMethodNotAllowed,
// 		},
// 		{
// 			name:           "invalid request body",
// 			method:         http.MethodPost,
// 			body:           "{bad json",
// 			mockSetup:      func(mock *mocks.MockGradeServiceI) {},
// 			expectedStatus: http.StatusBadRequest,
// 		},
// 		{
// 			name:   "service error",
// 			method: http.MethodPost,
// 			body: map[string]interface{}{
// 				"studentID": "1",
// 				"subjectID": "sub1",
// 				"semester":  1,
// 				"grade":     95,
// 			},
// 			mockSetup: func(mock *mocks.MockGradeServiceI) {
// 				mock.EXPECT().AddGrades("1", "sub1", 95, 1).Return(errors.New("db error"))
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var bodyBytes []byte
// 			switch v := tt.body.(type) {
// 			case string:
// 				bodyBytes = []byte(v)
// 			default:
// 				b, _ := json.Marshal(v)
// 				bodyBytes = b
// 			}

// 			req := httptest.NewRequest(tt.method, "/grades", bytes.NewReader(bodyBytes))
// 			w := httptest.NewRecorder()

// 			mockService := mocks.NewMockGradeServiceI(ctrl)
// 			tt.mockSetup(mockService)

// 			handler := handlers.NewGradeHandler(mockService)

// 			handler.AddGrade(w, req)

// 			if w.Code != tt.expectedStatus {
// 				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
// 			}
// 		})
// 	}
// }

// func TestGradeHandler_UpdateGrade(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	tests := []struct {
// 		name           string
// 		method         string
// 		body           any
// 		mockSetup      func(mock *mocks.MockGradeServiceI)
// 		expectedStatus int
// 	}{
// 		{
// 			name:   "success",
// 			method: http.MethodPatch,
// 			body: map[string]interface{}{
// 				"studentID": "1",
// 				"subjectID": "sub1",
// 				"new_grade": 90,
// 			},
// 			mockSetup: func(mock *mocks.MockGradeServiceI) {
// 				mock.EXPECT().UpdateGrade("1", "sub1", 90).Return(nil)
// 			},
// 			expectedStatus: http.StatusCreated,
// 		},
// 		{
// 			name:           "invalid method",
// 			method:         http.MethodGet,
// 			body:           nil,
// 			mockSetup:      func(mock *mocks.MockGradeServiceI) {},
// 			expectedStatus: http.StatusMethodNotAllowed,
// 		},
// 		{
// 			name:           "invalid request body",
// 			method:         http.MethodPatch,
// 			body:           "{bad json",
// 			mockSetup:      func(mock *mocks.MockGradeServiceI) {},
// 			expectedStatus: http.StatusBadRequest,
// 		},
// 		{
// 			name:   "service error",
// 			method: http.MethodPatch,
// 			body: map[string]interface{}{
// 				"studentID": "1",
// 				"subjectID": "sub1",
// 				"new_grade": 90,
// 			},
// 			mockSetup: func(mock *mocks.MockGradeServiceI) {
// 				mock.EXPECT().UpdateGrade("1", "sub1", 90).Return(errors.New("db error"))
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var bodyBytes []byte
// 			switch v := tt.body.(type) {
// 			case string:
// 				bodyBytes = []byte(v)
// 			default:
// 				b, _ := json.Marshal(v)
// 				bodyBytes = b
// 			}

// 			req := httptest.NewRequest(tt.method, "/grades", bytes.NewReader(bodyBytes))
// 			w := httptest.NewRecorder()

// 			mockService := mocks.NewMockGradeServiceI(ctrl)
// 			tt.mockSetup(mockService)

// 			handler := handlers.NewGradeHandler(mockService)
// 			handler.UpdateGrade(w, req)

// 			if w.Code != tt.expectedStatus {
// 				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
// 			}
// 		})
// 	}
// }
