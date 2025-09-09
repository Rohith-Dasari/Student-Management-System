package services

import "sms/models"

//go:generate mockgen -destination=../mocks/student_service_mock.go -package=mocks -source=student_service_interface.go
type StudentServiceI interface {
	CreateStudent(rollNumber, name, classID string, semester int) (*models.Students, error)
	UpdateStudent(studentID, name, rollnumber, classID string, semester int) error
}
