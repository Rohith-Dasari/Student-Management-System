package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sms/constants"
	"sms/handlers"
	"sms/mocks"
	gradeRepository "sms/repository/gradesRepository"
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
		role           constants.Role
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
func AddUserToContext(ctx context.Context, role constants.Role) context.Context {
	ctx = context.WithValue(ctx, constants.ContextUserRoleKey, role)
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
		role           constants.Role
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
		mockSetup      func()
		expectedStatus int
		role           constants.Role
		classID        string
		semester       string
	}{
		{
			name:     "successful retrieval",
			method:   http.MethodGet,
			classID:  "1",
			semester: "1",
			mockSetup: func() {
				mockGradeService.EXPECT().GetAverageOfClass("1", 1).Return(70.0, nil).Times(1)
			},
			expectedStatus: http.StatusOK,
			role:           "faculty",
		},
		{
			name:           "invalid method",
			method:         http.MethodPost,
			classID:        "1",
			semester:       "1",
			mockSetup:      func() {},
			expectedStatus: http.StatusMethodNotAllowed,
			role:           "faculty",
		},
		{
			name:     "invalid role",
			method:   http.MethodGet,
			classID:  "1",
			semester: "1",
			mockSetup: func() {
			},
			expectedStatus: http.StatusForbidden,
			role:           "admin",
		},
		{
			name:           "invalid classID",
			method:         http.MethodGet,
			classID:        "",
			semester:       "1",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			role:           "faculty",
		},
		{
			name:           "invalid semester",
			method:         http.MethodGet,
			classID:        "1",
			semester:       "invalid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			role:           "faculty",
		},
		{
			name:     "service error",
			method:   http.MethodGet,
			classID:  "1",
			semester: "1",
			mockSetup: func() {
				mockGradeService.EXPECT().GetAverageOfClass("1", 1).Return(0.0, errors.New("service error")).Times(1)
			},
			expectedStatus: http.StatusBadRequest,
			role:           "faculty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, fmt.Sprintf("/classes/%s/semesters/%s/average", tt.classID, tt.semester), nil)
			req.SetPathValue("classID", tt.classID)
			req.SetPathValue("semester", tt.semester)

			req = req.WithContext(AddUserToContext(req.Context(), tt.role))
			tt.mockSetup()

			rr := httptest.NewRecorder()

			handler.GetAverageOfClass(rr, req)

			if rr.Code != tt.expectedStatus {
				var resp map[string]any
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				if err != nil {
					t.Errorf("unmarshalling response body failed: %v", err)
				}
				t.Errorf("Test failed: %s. Expected status %d, got %d. Response: %v", tt.name, tt.expectedStatus, rr.Code, resp)
			}
		})
	}
}

func TestHandler_GetToppers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGradeService := mocks.NewMockGradeServiceI(ctrl)

	handler := handlers.NewGradeHandler(mockGradeService)
	tests := []struct {
		name           string
		classID        string
		semester       string
		topLimit       string
		mockSetup      func()
		expectedStatus int
		role           constants.Role
	}{
		{
			name:     "successful retrieval",
			classID:  "1",
			semester: "1",
			topLimit: "3",
			mockSetup: func() {
				mockGradeService.EXPECT().GetToppers("1", 1, 3).Return([]gradeRepository.StudentAverage{}, nil).Times(1)
			},
			expectedStatus: http.StatusOK,
			role:           "faculty",
		},
		{
			name:     "invalid role",
			classID:  "1",
			semester: "1",
			topLimit: "3",
			mockSetup: func() {

			},
			expectedStatus: http.StatusForbidden,
			role:           "admin",
		},
		{
			name:     "invalid classID",
			classID:  "",
			semester: "1",
			topLimit: "3",
			mockSetup: func() {
			},
			expectedStatus: http.StatusBadRequest,
			role:           "faculty",
		},
		{
			name:     "invalid semester",
			classID:  "1",
			semester: "invalid",
			topLimit: "3",
			mockSetup: func() {
			},
			expectedStatus: http.StatusBadRequest,
			role:           "faculty",
		},
		{
			name:     "invalid top limit",
			classID:  "1",
			semester: "1",
			topLimit: "invalid",
			mockSetup: func() {
			},
			expectedStatus: http.StatusBadRequest,
			role:           "faculty",
		},
		{
			name:     "service error",
			classID:  "1",
			semester: "1",
			topLimit: "3",
			mockSetup: func() {
				mockGradeService.EXPECT().GetToppers("1", 1, 3).Return(nil, errors.New("service error")).Times(1)
			},
			expectedStatus: http.StatusBadRequest,
			role:           "faculty",
		},
		{
			name:           "missing top param",
			classID:        "1",
			semester:       "1",
			topLimit:       "",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			role:           "faculty",
		},
		{
			name:           "missing user role in context",
			classID:        "1",
			semester:       "1",
			topLimit:       "3",
			mockSetup:      func() {},
			expectedStatus: http.StatusForbidden,
			role:           "",
		},
		{
			name:           "negative top limit",
			classID:        "1",
			semester:       "1",
			topLimit:       "-5",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			role:           "faculty",
		},
		{
			name:           "negative semester",
			classID:        "1",
			semester:       "-1",
			topLimit:       "3",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			role:           "faculty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/classes/%s/semesters/%s/toppers?top=%s", tt.classID, tt.semester, tt.topLimit), nil)

			req.SetPathValue("classID", tt.classID)
			req.SetPathValue("semester", tt.semester)

			req = req.WithContext(AddUserToContext(req.Context(), tt.role))
			tt.mockSetup()

			rr := httptest.NewRecorder()

			handler.GetToppers(rr, req)

			if rr.Code != tt.expectedStatus {
				var resp map[string]any
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				if err != nil {
					t.Errorf("unmarshalling response body failed: %v", err)
				}
				t.Errorf("Test failed: %s. Expected status %d, got %d. Response: %v", tt.name, tt.expectedStatus, rr.Code, resp)
			}
		})
	}
}
