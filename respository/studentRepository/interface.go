package studentsRepository

import "sms/models"

type StudentRepositoryI interface {
	GetAllStudentsOfClass(classID string, semester int) ([]models.Students, error)
	AddStudent(uuid, rollNumber, name, classID string, semester int) error
	UpdateStudent(studentID, name, rollnumber, classID string, semester int) error
	GetStudentByID(studentID string) (*models.Students, error)
	GetStudentByRollNumber(rollNumber string) (*models.Students, error)
}
