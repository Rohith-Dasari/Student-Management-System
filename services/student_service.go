package services

import (
	"errors"
	"fmt"
	"sms/models"
	studentRepo "sms/repository/studentRepository"

	"github.com/google/uuid"
)

type StudentService struct {
	sr studentRepo.StudentRepositoryI
}

func NewStudentService(sr studentRepo.StudentRepositoryI) StudentService {
	return StudentService{sr}
}

func (ss *StudentService) CreateStudent(rollNumber, name, classID string, semester int) (*models.Students, error) {
	//rollNumber check
	student, _ := ss.sr.GetStudentByRollNumber(rollNumber)
	if student != nil {
		return nil, fmt.Errorf("student %v already exists with given rollnumber", student)
	}
	//name not empty
	if name == "" {
		return nil, errors.New("name can't be empty")
	}
	uuid := uuid.New().String()
	err := ss.sr.AddStudent(uuid, rollNumber, name, classID, semester)
	if err != nil {
		return nil, err
	}
	newStudent := models.Students{
		StudentID:  uuid,
		RollNumber: rollNumber,
		ClassID:    classID,
		Semester:   semester,
		Name:       name,
	}
	return &newStudent, nil
}

func (ss *StudentService) UpdateStudent(studentID, name, rollnumber, classID string, semester int) error {
	student, _ := ss.sr.GetStudentByID(studentID)

	if name != "" {
		student.Name = name
	}
	if rollnumber != "" {
		student.RollNumber = rollnumber
	}
	if classID != "" {
		student.ClassID = classID
	}
	if semester != 0 {
		student.Semester = semester
	}

	err := ss.sr.UpdateStudent(studentID, student.Name, student.RollNumber, student.ClassID, student.Semester)
	return err
}
