package services

import (
	"errors"
	gradeRepository "sms/repository/gradesRepository"
)

type GradeService struct {
	gr gradeRepository.GradeRepositoryI
}

func NewGradeService(gr gradeRepository.GradeRepositoryI) *GradeService {
	return &GradeService{gr}
}

func (gs *GradeService) GetAverageOfClass(classID string, semester int) (float64, error) {
	averageOfClass, err := gs.gr.GetClassAverage(classID, semester)
	if err != nil {
		return 0, err
	}
	return averageOfClass, nil
}

func (gs *GradeService) GetToppers(classID string, semester int, top int) ([]gradeRepository.StudentAverage, error) {
	toppers, err := gs.gr.GetToppers(classID, semester, top)
	if err != nil {
		return nil, err
	}
	return toppers, nil
}

func (gs *GradeService) AddGrades(studentID string, subjectID string, grade int, semester int) error {
	if grade < 0 {
		return errors.New("grade can't be negative")
	}
	err := gs.gr.AddGrades(studentID, subjectID, grade, semester)
	return err
}

func (gs *GradeService) UpdateGrade(studentID string, subjectID string, newGrade int) error {
	if newGrade < 0 {
		return errors.New("grade can't be negative")
	}
	err := gs.gr.UpdateGrade(studentID, subjectID, newGrade)
	return err
}
