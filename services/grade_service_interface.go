package services

import (
	gradeRepository "sms/repository/gradesRepository"
)

//go:generate mockgen -destination=../mocks/grade_service_mock.go -package=mocks -source=grade_service_interface.go
type GradeServiceI interface {
	GetAverageOfClass(classID string, semester int) (float64, error)
	GetToppers(classID string, semester int, top int) ([]gradeRepository.StudentAverage, error)
	AddGrades(studentID string, subjectID string, Grade int, semester int) error
	UpdateGrade(studentID string, subjectID string, newGrade int) error
}
