package studentsRepository

import (
	"database/sql"
	"sms/models"
)

type StudentRepo struct {
	db *sql.DB
}

func NewStudentRepo(db *sql.DB) *StudentRepo {
	return &StudentRepo{db}
}

func (sr *StudentRepo) AddStudent(uuid, rollNumber, name, classID string, semester int) error {
	_, err := sr.db.Exec(`insert into students values(?,?,?,?,?)`, uuid, name, rollNumber, classID, semester)
	return err
}

func (sr *StudentRepo) UpdateStudent(studentID, name, rollnumber, classID string, semester int) error {
	stmt := `update students set Name=?,RollNumber=?,ClassID=?,semester=? where StudentID=?`
	_, err := sr.db.Exec(stmt, name, rollnumber, classID, semester, studentID)
	return err
}

func (sr *StudentRepo) GetStudentByID(studentID string) (*models.Students, error) {
	stmt := `select StudentID,Name,RollNumber,ClassID,semester from students where StudentID=?`
	row := sr.db.QueryRow(stmt, studentID)
	var s models.Students
	err := row.Scan(&s.StudentID, &s.Name, &s.RollNumber, &s.ClassID, &s.Semester)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}
func (sr *StudentRepo) GetStudentByRollNumber(rollNumber string) (*models.Students, error) {
	stmt := `select StudentID,Name,RollNumber,ClassID,semester from students where RollNumber=?`
	row := sr.db.QueryRow(stmt, rollNumber)
	var s models.Students
	err := row.Scan(&s.StudentID, &s.Name, &s.RollNumber, &s.ClassID, &s.Semester)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}
