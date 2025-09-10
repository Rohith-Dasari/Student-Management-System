package services_test

import (
	"errors"
	"testing"

	"sms/mocks"
	mockrepo "sms/mocks"
	"sms/models"
	"sms/services"

	"github.com/golang/mock/gomock"
)

func TestGetAverageGradeOfStudent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGradeRepo := mockrepo.NewMockGradeRepositoryI(ctrl)
	mockStudentRepo := mockrepo.NewMockStudentRepositoryI(ctrl)

	gs := services.NewGradeService(mockGradeRepo, mockStudentRepo)

	studentID := "stu1"
	semester := 1
	mockGradeRepo.EXPECT().GetSemesterGrades(studentID, semester).Return([]int{80, 90, 100}, nil)

	avg, err := gs.GetAverageGradeOfStudent(studentID, semester)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if avg != 90 {
		t.Errorf("expected average 90, got %d", avg)
	}
}

func TestGetAverageGradesOfEachStudentOfClass(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGradeRepo := mockrepo.NewMockGradeRepositoryI(ctrl)
	mockStudentRepo := mockrepo.NewMockStudentRepositoryI(ctrl)

	gs := services.NewGradeService(mockGradeRepo, mockStudentRepo)

	classID := "class1"
	semester := 1

	students := []models.Students{
		{StudentID: "stu1", Name: "A"},
		{StudentID: "stu2", Name: "B"},
	}

	mockStudentRepo.EXPECT().GetAllStudentsOfClass(classID, semester).Return(students, nil)
	mockGradeRepo.EXPECT().GetSemesterGrades("stu1", semester).Return([]int{80, 90}, nil)
	mockGradeRepo.EXPECT().GetSemesterGrades("stu2", semester).Return([]int{70, 100}, nil)

	result, err := gs.GetAverageGradesOfEachStudentOfClass(classID, semester)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 students, got %d", len(result))
	}

	if result["stu1"] != 85 || result["stu2"] != 85 {
		t.Errorf("expected averages 85, got %v", result)
	}
}

func TestGetTopThree(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGradeRepo := mockrepo.NewMockGradeRepositoryI(ctrl)
	mockStudentRepo := mockrepo.NewMockStudentRepositoryI(ctrl)

	gs := services.NewGradeService(mockGradeRepo, mockStudentRepo)

	classID := "class1"
	semester := 1

	students := []models.Students{
		{StudentID: "s1", Name: "Alice"},
		{StudentID: "s2", Name: "Bob"},
		{StudentID: "s3", Name: "Charlie"},
		{StudentID: "s4", Name: "David"},
	}

	mockStudentRepo.EXPECT().GetAllStudentsOfClass(classID, semester).Return(students, nil)
	mockGradeRepo.EXPECT().GetSemesterGrades("s1", semester).Return([]int{90}, nil)
	mockGradeRepo.EXPECT().GetSemesterGrades("s2", semester).Return([]int{95}, nil)
	mockGradeRepo.EXPECT().GetSemesterGrades("s3", semester).Return([]int{85}, nil)
	mockGradeRepo.EXPECT().GetSemesterGrades("s4", semester).Return([]int{80}, nil)

	mockStudentRepo.EXPECT().GetStudentByID("s2").Return(&students[1], nil)
	mockStudentRepo.EXPECT().GetStudentByID("s1").Return(&students[0], nil)
	mockStudentRepo.EXPECT().GetStudentByID("s3").Return(&students[2], nil)

	topThree, err := gs.GetTopThree(classID, semester)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(topThree) != 3 {
		t.Errorf("expected 3 students, got %d", len(topThree))
	}

	if topThree[0].StudentID != "s2" || topThree[1].StudentID != "s1" || topThree[2].StudentID != "s3" {
		t.Errorf("top three order is wrong: %v", topThree)
	}
}

func TestAddGrades_And_UpdateGrade(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGradeRepo := mockrepo.NewMockGradeRepositoryI(ctrl)
	mockStudentRepo := mockrepo.NewMockStudentRepositoryI(ctrl)

	gs := services.NewGradeService(mockGradeRepo, mockStudentRepo)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGrades := mocks.NewMockGradeRepositoryI(ctrl)
	mockStudents := mocks.NewMockStudentRepositoryI(ctrl)

	service := services.NewGradeService(mockGrades, mockStudents)

	classID := "class1"
	semester := 1

	t.Run("success - average calculated", func(t *testing.T) {
		students := []models.Students{
			{StudentID: "s1", Name: "Alice"},
			{StudentID: "s2", Name: "Bob"},
		}

		mockStudents.EXPECT().
			GetAllStudentsOfClass(classID, semester).
			Return(students, nil)

		mockGrades.EXPECT().
			GetSemesterGrades("s1", semester).
			Return([]int{80, 90}, nil)

		mockGrades.EXPECT().
			GetSemesterGrades("s2", semester).
			Return([]int{70, 100}, nil)

		avg, err := service.GetAverageOfClass(classID, semester)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		expected := 85
		if avg != expected {
			t.Errorf("expected %d, got %d", expected, avg)
		}
	})

	t.Run("error - student repo fails", func(t *testing.T) {
		mockStudents.EXPECT().
			GetAllStudentsOfClass(classID, semester).
			Return(nil, errors.New("db error"))

		avg, err := service.GetAverageOfClass(classID, semester)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if avg != 0 {
			t.Errorf("expected 0, got %d", avg)
		}
	})

	t.Run("error - grade repo fails", func(t *testing.T) {
		students := []models.Students{
			{StudentID: "s1", Name: "Alice"},
		}

		mockStudents.EXPECT().
			GetAllStudentsOfClass(classID, semester).
			Return(students, nil)

		mockGrades.EXPECT().
			GetSemesterGrades("s1", semester).
			Return(nil, errors.New("db error"))

		avg, err := service.GetAverageOfClass(classID, semester)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if avg != 0 {
			t.Errorf("expected 0, got %d", avg)
		}
	})
}
