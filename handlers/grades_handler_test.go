package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sms/handlers"
	"sms/middleware"
	"sms/mocks"
	"sms/services"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestGradeHandler_AddGrade(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGradeService := mocks.NewMockGradeServiceI(ctrl)

	handler := handlers.NewGradeHandler(mockGradeService)
	tests := []struct {
		name           string
		method         string
		body           any
		role           string
		mockService    func()
		expectedStatus int
	}{
		{
			name:   "faculty adds grade successfully",
			method: http.MethodPost,
			body: map[string]any{
				"studentID": "1",
				"subjectID": "sub1",
				"semester":  1,
				"grade":     95,
			},
			role: "faculty",
			mockService: func() {
				mockGradeService.EXPECT().AddGrades("1", "sub1", 95, 1).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:   "admin tries to add grade unsuccessfully",
			method: http.MethodPost,
			body: map[string]any{
				"studentID": "1",
				"subjectID": "sub2",
				"semester":  1,
				"grade":     95,
			},
			role:           "admin",
			mockService:    func() {},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:   "service error",
			method: http.MethodPost,
			body: map[string]any{
				"studentID": "1",
				"subjectID": "sub1",
				"semester":  1,
				"grade":     95,
			},
			role: "faculty",
			mockService: func() {
				mockGradeService.EXPECT().AddGrades("1", "sub1", 95, 1).Return(errors.New("grade already exists"))
			},
			expectedStatus: http.StatusBadRequest,
		}, {
			name:   "wrong method",
			method: http.MethodGet,
			body: map[string]any{
				"studentID": "1",
				"subjectID": "sub1",
				"semester":  1,
				"grade":     95,
			},
			role:           "faculty",
			mockService:    func() {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "invalid request body",
			method: http.MethodPost,
			body: map[any]any{
				"studentID": 1789,
				"grade":     "something",
			},
			role:           "faculty",
			mockService:    func() {},
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
			req := httptest.NewRequest(tt.method, "/grades", bytes.NewReader(reqBody))

			req = req.WithContext(AddUserToContext(req.Context(), tt.role))

			rr := httptest.NewRecorder()
			tt.mockService()
			handler.AddGrade(rr, req)
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

		})

	}

}
func AddUserToContext(ctx context.Context, role string) context.Context {
	ctx = context.WithValue(ctx, middleware.ContextUserRoleKey, role)
	return ctx
}

func TestGradeHandler_UpdateGrade(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGradeService := mocks.NewMockGradeServiceI(ctrl)

	handler := handlers.NewGradeHandler(mockGradeService)

	tests := []struct {
		name           string
		method         string
		body           any
		mockSetup      func()
		expectedStatus int
		role           string
	}{
		{
			name:   "successful update",
			method: http.MethodPatch,
			body: map[string]any{
				"studentID": "1",
				"subjectID": "sub1",
				"new_grade": 95,
			},
			mockSetup: func() {
				mockGradeService.EXPECT().UpdateGrade("1", "sub1", 95).Return(nil)
			},
			expectedStatus: http.StatusOK,
			role:           "faculty",
		},
		{
			name:   "invalid method",
			method: http.MethodGet,
			body: map[string]any{
				"studentID": "1",
				"subjectID": "sub1",
				"new_grade": 95,
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusMethodNotAllowed,
			role:           "faculty",
		},
		{
			name:   "invalid role",
			method: http.MethodPatch,
			body: map[string]any{
				"studentID": "1",
				"subjectID": "sub1",
				"new_grade": 95,
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusForbidden,
			role:           "admin",
		},
		{
			name:   "invalid request body",
			method: http.MethodPatch,
			body: map[string]any{
				"studentID": 1,
				"subjectID": "sub1",
				"new_grade": "invalid",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			role:           "faculty",
		},
		{
			name:   "service error",
			method: http.MethodPatch,
			body: map[string]any{
				"studentID": "invalid",
				"subjectID": "sub1",
				"new_grade": 70,
			},
			mockSetup: func() {
				mockGradeService.EXPECT().UpdateGrade("invalid", "sub1", 70).Return(errors.New("invalid studentID"))
			},
			expectedStatus: http.StatusBadRequest,
			role:           "faculty",
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
			req := httptest.NewRequest(tt.method, "/grades", bytes.NewReader(reqBody))

			req = req.WithContext(AddUserToContext(req.Context(), tt.role))

			rr := httptest.NewRecorder()
			tt.mockSetup()
			handler.UpdateGrade(rr, req)
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

		})

	}

}

func TestHandler_GetAverageOfClass(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGradeService := mocks.NewMockGradeServiceI(ctrl)

	handler := handlers.NewGradeHandler(mockGradeService)
	tests := []struct {
		name           string
		method         string
		query          string
		mockSetup      func()
		expectedStatus int
		role           string
	}{
		{
			name:   "successful retrieval",
			method: http.MethodGet,
			query:  "?classID=1&semester=1",
			mockSetup: func() {
				mockGradeService.EXPECT().GetAverageOfClass(gomock.Any(), gomock.Any()).Return(70, nil)
			},
			expectedStatus: http.StatusOK,
			role:           "faculty",
		},
		{
			name:           "invalid method",
			method:         http.MethodPost,
			query:          "?classID=1&semester=1",
			mockSetup:      func() {},
			expectedStatus: http.StatusMethodNotAllowed,
			role:           "faculty",
		},
		{
			name:   "invalid role",
			method: http.MethodGet,
			query:  "?classID=1&semester=1",
			mockSetup: func() {
			},
			expectedStatus: http.StatusForbidden,
			role:           "admin",
		},
		{
			name:           "invalid classID",
			method:         http.MethodGet,
			query:          "?semester=1",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			role:           "faculty",
		},
		{
			name:           "invalid semester",
			method:         http.MethodGet,
			query:          "?classID=1&semester=invalid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			role:           "faculty",
		},
		{
			name:   "service error",
			method: http.MethodGet,
			query:  "?classID=1&semester=7",
			mockSetup: func() {
				mockGradeService.EXPECT().GetAverageOfClass(gomock.Any(), gomock.Any()).Return(0, errors.New("service error"))
			},
			expectedStatus: http.StatusBadRequest,
			role:           "faculty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(tt.method, "/grades"+tt.query, nil)

			req = req.WithContext(AddUserToContext(req.Context(), tt.role))
			tt.mockSetup()

			rr := httptest.NewRecorder()
			handler.GetAverageOfClass(rr, req)
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestHandler_GetTopThree(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGradeService := mocks.NewMockGradeServiceI(ctrl)

	handler := handlers.NewGradeHandler(mockGradeService)
	tests := []struct {
		name           string
		query          string
		mockSetup      func()
		expectedStatus int
		role           string
	}{{
		name:  "successful retrieval",
		query: "?classID=1&semester=1",
		mockSetup: func() {
			mockGradeService.EXPECT().GetTopThree(gomock.Any(), gomock.Any()).Return([]services.GradeReponse{}, nil)
		},
		expectedStatus: http.StatusOK,
		role:           "faculty",
	},
		{
			name:  "invalid role",
			query: "?classID=1&semester=1",
			mockSetup: func() {

			},
			expectedStatus: http.StatusForbidden,
			role:           "admin",
		},
		{
			name:  "invalid classID",
			query: "?classID=&semester=1",
			mockSetup: func() {
			},
			expectedStatus: http.StatusBadRequest,
			role:           "faculty",
		},
		{
			name:  "invalid semster",
			query: "?classID=1&semester=invalid",
			mockSetup: func() {
			},
			expectedStatus: http.StatusBadRequest,
			role:           "faculty",
		},
		{
			name:  "service error",
			query: "?classID=1&semester=2",
			mockSetup: func() {
				mockGradeService.EXPECT().GetTopThree(gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusBadRequest,
			role:           "faculty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, "/grades/toppers"+tt.query, nil)

			req = req.WithContext(AddUserToContext(req.Context(), tt.role))
			tt.mockSetup()

			rr := httptest.NewRecorder()
			handler.GetTopThree(rr, req)
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})

	}
}
