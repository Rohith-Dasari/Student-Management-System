package gradeRepository

//go:generate mockgen -destination=../../mocks/grade_repo_mock.go -package=mocks -source=interface.go
type GradeRepositoryI interface {
	GetSemesterGrades(studentID string, semester int) ([]int, error)
	AddGrades(studentID string, subjectID string, Grade int, semester int) error
	UpdateGrade(studentID string, subjectID string, newGrade int) error
}
