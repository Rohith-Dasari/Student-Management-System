package services_test

import (
	"errors"
	"reflect"
	"testing"

	"sms/mocks"
	mockrepo "sms/mocks"
	gradeRepository "sms/repository/gradesRepository"
	"sms/services"

	"go.uber.org/mock/gomock"
)

func TestGetToppers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockGradeRepositoryI(ctrl)
	gradeService := services.NewGradeService(mockRepo)
	expectedToppers := []gradeRepository.StudentAverage{
		{StudentID: "student1", Average: 95.5},
		{StudentID: "student2", Average: 92.0},
		{StudentID: "student3", Average: 88.5},
	}

	tests := []struct {
		name            string
		classID         string
		semester        int
		top             int
		mockSetup       func()
		expectedToppers []gradeRepository.StudentAverage
		expectedError   error
	}{
		{
			name:     "Successful retrieval of toppers",
			classID:  "CS101",
			semester: 1,
			top:      3,
			mockSetup: func() {
				mockRepo.EXPECT().GetToppers("CS101", 1, 3).Return(expectedToppers, nil).Times(1)
			},
			expectedToppers: expectedToppers,
			expectedError:   nil,
		},
		{
			name:     "Repository returns an error",
			classID:  "CS102",
			semester: 2,
			top:      5,
			mockSetup: func() {
				mockRepo.EXPECT().GetToppers("CS102", 2, 5).Return(nil, errors.New("database error")).Times(1)
			},
			expectedToppers: nil,
			expectedError:   errors.New("database error"),
		},
		{
			name:     "Empty result from repository",
			classID:  "CS103",
			semester: 3,
			top:      10,
			mockSetup: func() {
				mockRepo.EXPECT().GetToppers("CS103", 3, 10).Return([]gradeRepository.StudentAverage{}, nil).Times(1)
			},
			expectedToppers: []gradeRepository.StudentAverage{},
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			toppers, err := gradeService.GetToppers(tt.classID, tt.semester, tt.top)

			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(toppers, tt.expectedToppers) {
					t.Errorf("expected toppers: %v, got: %v", tt.expectedToppers, toppers)
				}
			}
		})
	}
}

func TestAddGrades_And_UpdateGrade(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGradeRepo := mockrepo.NewMockGradeRepositoryI(ctrl)

	gs := services.NewGradeService(mockGradeRepo)

	mockGradeRepo.EXPECT().AddGrades("s1", "sub1", 90, 1).Return(nil)
	if err := gs.AddGrades("s1", "sub1", 90, 1); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := gs.AddGrades("s1", "sub1", -5, 1); err == nil {
		t.Errorf("expected error for negative grade")
	}

	mockGradeRepo.EXPECT().UpdateGrade("s1", "sub1", 95).Return(nil)
	if err := gs.UpdateGrade("s1", "sub1", 95); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := gs.UpdateGrade("s1", "sub1", -10); err == nil {
		t.Errorf("expected error for negative grade")
	}
}

func TestGetAverageOfClass(t *testing.T) {
	tests := []struct {
		name          string
		classID       string
		semester      int
		mockSetup     func(mockRepo *mocks.MockGradeRepositoryI)
		expectedAvg   float64
		expectedError error
	}{
		{
			name:     "Successful retrieval of average",
			classID:  "CS101",
			semester: 1,
			mockSetup: func(mockRepo *mocks.MockGradeRepositoryI) {
				mockRepo.EXPECT().GetClassAverage("CS101", 1).Return(85.5, nil).Times(1)
			},
			expectedAvg:   85.5,
			expectedError: nil,
		},
		{
			name:     "Repository returns an error",
			classID:  "CS102",
			semester: 2,
			mockSetup: func(mockRepo *mocks.MockGradeRepositoryI) {
				mockRepo.EXPECT().GetClassAverage("CS102", 2).Return(0.0, errors.New("database error")).Times(1)
			},
			expectedAvg:   0.0,
			expectedError: errors.New("database error"),
		},
		{
			name:     "No class ID provided",
			classID:  "",
			semester: 1,
			mockSetup: func(mockRepo *mocks.MockGradeRepositoryI) {
				mockRepo.EXPECT().GetClassAverage("", 1).Return(0.0, errors.New("invalid classID from repo")).Times(1)
			},
			expectedAvg:   0.0,
			expectedError: errors.New("invalid classID from repo"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockGradeRepositoryI(ctrl)
			gradeService := services.NewGradeService(mockRepo)

			tt.mockSetup(mockRepo)

			average, err := gradeService.GetAverageOfClass(tt.classID, tt.semester)

			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if average != tt.expectedAvg {
					t.Errorf("expected average: %f, got: %f", tt.expectedAvg, average)
				}
			}
		})
	}
}
