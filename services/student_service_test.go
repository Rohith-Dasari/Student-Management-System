package services_test

import (
	"errors"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"

	mockrepo "sms/mocks"
	"sms/models"
	"sms/services"
)

func TestCreateStudent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockStudentRepositoryI(ctrl)
	svc := services.NewStudentService(mockRepo)

	mockRepo.EXPECT().GetStudentByRollNumber("101").Return(nil, nil)
	mockRepo.EXPECT().AddStudent(gomock.Any(), "101", "Rohith", "CSE", 5).Return(nil)

	student, err := svc.CreateStudent("101", "Rohith", "CSE", 5)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if student == nil {
		t.Fatal("expected student, got nil")
	}
	if student.RollNumber != "101" {
		t.Errorf("expected rollnumber 101, got %v", student.RollNumber)
	}
	if student.Name != "Rohith" {
		t.Errorf("expected name Rohith, got %v", student.Name)
	}
}

func TestCreateStudent_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockStudentRepositoryI(ctrl)
	svc := services.NewStudentService(mockRepo)

	existing := &models.Students{StudentID: "123", RollNumber: "101", Name: "Old"}

	mockRepo.EXPECT().GetStudentByRollNumber("101").Return(existing, nil)

	student, err := svc.CreateStudent("101", "Rohith", "CSE", 5)

	if student != nil {
		t.Fatal("expected nil student, got value")
	}
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected already exists error, got %v", err)
	}
}

func TestCreateStudent_NameEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockStudentRepositoryI(ctrl)
	svc := services.NewStudentService(mockRepo)

	mockRepo.EXPECT().GetStudentByRollNumber("102").Return(nil, nil)

	student, err := svc.CreateStudent("102", "", "CSE", 5)

	if student != nil {
		t.Fatal("expected nil student, got value")
	}
	if err == nil || err.Error() != "name can't be empty" {
		t.Fatalf("expected name can't be empty error, got %v", err)
	}
}

func TestCreateStudent_AddStudentError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockStudentRepositoryI(ctrl)
	svc := services.NewStudentService(mockRepo)

	mockRepo.EXPECT().GetStudentByRollNumber("103").Return(nil, nil)
	mockRepo.EXPECT().AddStudent(gomock.Any(), "103", "New", "CSE", 5).Return(errors.New("db error"))

	student, err := svc.CreateStudent("103", "New", "CSE", 5)

	if student != nil {
		t.Fatal("expected nil student, got value")
	}
	if err == nil || err.Error() != "db error" {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestUpdateStudent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockStudentRepositoryI(ctrl)
	svc := services.NewStudentService(mockRepo)

	existing := &models.Students{StudentID: "123", RollNumber: "101", Name: "Old", ClassID: "CSE", Semester: 5}

	mockRepo.EXPECT().GetStudentByID("123").Return(existing, nil)
	mockRepo.EXPECT().UpdateStudent("123", "NewName", "101", "CSE", 5).Return(nil)

	err := svc.UpdateStudent("123", "NewName", "", "", 0)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUpdateStudent_UpdateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockStudentRepositoryI(ctrl)
	svc := services.NewStudentService(mockRepo)

	existing := &models.Students{StudentID: "123", RollNumber: "101", Name: "Old", ClassID: "CSE", Semester: 5}

	mockRepo.EXPECT().GetStudentByID("123").Return(existing, nil)
	mockRepo.EXPECT().UpdateStudent("123", "FailName", "101", "CSE", 5).Return(errors.New("update failed"))

	err := svc.UpdateStudent("123", "FailName", "", "", 0)

	if err == nil || err.Error() != "update failed" {
		t.Fatalf("expected update failed error, got %v", err)
	}
}
