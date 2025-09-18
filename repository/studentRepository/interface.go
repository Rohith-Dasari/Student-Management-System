package studentsRepository

import "sms/models"

//go:generate mockgen -destination=../../mocks/student_repo_mock.go -package=mocks -source=interface.go
type StudentRepositoryI interface {
	AddStudent(uuid, rollNumber, name, classID string, semester int) error
	UpdateStudent(studentID, name, rollnumber, classID string, semester int) error
	GetStudentByID(studentID string) (*models.Students, error)
	GetStudentByRollNumber(rollNumber string) (*models.Students, error)
}
