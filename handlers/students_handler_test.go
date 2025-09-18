package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sms/constants"
	"sms/handlers"
	"sms/mocks"
	"sms/models"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestStudentHandler_AddStudent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStudentService := mocks.NewMockStudentServiceI(ctrl)

	handler := handlers.NewStudentHandler(mockStudentService)

	tests := []struct {
		name           string
		body           any
		role           constants.Role
		mockService    func()
		expectedStatus int
	}{
		{
			name: "admin adds student successfully",
			body: map[string]any{
				"roll_number": "1",
				"name":        "rohith",
				"classID":     "1",
				"semester":    7,
			},
			role: "admin",
			mockService: func() {
				mockStudentService.EXPECT().CreateStudent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.Students{}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "faculty can't add student",
			body: map[string]any{
				"roll_number": "1",
				"name":        "rohith",
				"classID":     "1",
				"semester":    7,
			},
			role: "faculty",
			mockService: func() {
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "service error",
			body: map[string]any{
				"roll_number": "1",
				"name":        "rohith",
				"classID":     "1",
				"semester":    7,
			},
			role: "admin",
			mockService: func() {
				mockStudentService.EXPECT().CreateStudent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid request body",
			body: map[string]any{
				"roll_number": 123,
				"name":        nil,
				"classID":     1,
				"semester":    "invalid",
			},
			role: "admin",
			mockService: func() {
				// mockStudentService.EXPECT().CreateStudent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Students{}, nil)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			if s, ok := tt.body.(string); ok {
				reqBody = []byte(s)
			} else {
				reqBody, _ = json.Marshal(tt.body)
			}
			req := httptest.NewRequest(http.MethodPost, "/students", bytes.NewReader(reqBody))

			req = req.WithContext(AddUserToContext(req.Context(), tt.role))

			rr := httptest.NewRecorder()
			tt.mockService()
			handler.AddStudent(rr, req)
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

		})
	}
}

func TestStudentHandler_UpdateStudent(t *testing.T) {
	tests := []struct {
		name           string
		body           any
		expectedStatus int
		mockService    func(*mocks.MockStudentServiceI)
		role           constants.Role
		studentID      string
	}{
		{
			name:           "successful update",
			role:           "admin",
			expectedStatus: http.StatusOK,
			body: map[string]any{
				"roll_number": "1",
				"name":        "rohith",
				"classID":     "1",
				"semester":    7,
			},
			mockService: func(mockStudentService *mocks.MockStudentServiceI) {
				mockStudentService.EXPECT().UpdateStudent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			studentID: "1",
		},
		{
			name:           "invalid role",
			role:           "faculty",
			expectedStatus: http.StatusForbidden,
			body: map[string]any{
				"roll_number": "1",
				"name":        "rohith",
				"classID":     "1",
				"semester":    7,
			},
			mockService: func(mockStudentService *mocks.MockStudentServiceI) {
			},
			studentID: "1",
		},
		{
			name:           "invalid studentID",
			role:           "admin",
			expectedStatus: http.StatusBadRequest,
			body: map[string]any{
				"roll_number": "1",
				"name":        "rohith",
				"classID":     "1",
				"semester":    7,
			},
			mockService: func(mockStudentService *mocks.MockStudentServiceI) {
			},
			studentID: "",
		},
		{
			name:           "invalid request body",
			role:           "admin",
			expectedStatus: http.StatusBadRequest,
			body: map[string]any{
				"roll_number": 1,
				"classID":     7,
				"semester":    "invalid",
			},
			mockService: func(mockStudentService *mocks.MockStudentServiceI) {
			},
			studentID: "1",
		},
		{
			name:           "service error",
			role:           "admin",
			expectedStatus: http.StatusBadRequest,
			body: map[string]any{
				"roll_number": "1",
				"name":        "rohith",
				"classID":     "1",
				"semester":    7,
			},
			mockService: func(mockStudentService *mocks.MockStudentServiceI) {
				mockStudentService.EXPECT().
					UpdateStudent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("service error"))
			},
			studentID: "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStudentService := mocks.NewMockStudentServiceI(ctrl)
			handler := handlers.NewStudentHandler(mockStudentService)

			var reqBody []byte
			if s, ok := tt.body.(string); ok {
				reqBody = []byte(s)
			} else {
				reqBody, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(http.MethodPatch, "/students/"+tt.studentID, bytes.NewReader(reqBody))
			req = req.WithContext(AddUserToContext(req.Context(), tt.role))

			req.SetPathValue("studentID", tt.studentID)

			tt.mockService(mockStudentService)
			rr := httptest.NewRecorder()

			handler.UpdateStudent(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
