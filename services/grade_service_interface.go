package services

//go:generate mockgen -destination=../mocks/grade_service_mock.go -package=mocks -source=grade_service_interface.go
type GradeServiceI interface {
	GetAverageOfClass(classID string, semester int) (int, error)
	GetTopThree(classID string, semester int) ([]GradeReponse, error)
	AddGrades(studentID string, subjectID string, Grade int, semester int) error
	UpdateGrade(studentID string, subjectID string, newGrade int) error
}
