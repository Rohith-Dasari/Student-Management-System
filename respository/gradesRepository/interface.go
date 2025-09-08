package gradeRepository

type GradeRepositoryI interface {
	GetSemesterGrades(studentID string, semester int) ([]int, error)
	AddGrades(studentID string, subjectID string, Grade int, semester int) error
}
