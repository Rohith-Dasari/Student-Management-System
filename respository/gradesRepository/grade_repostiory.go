package gradeRepository

import "database/sql"

type GradeRepo struct {
	db *sql.DB
}

func NewGradeRepo(db *sql.DB) *GradeRepo {
	return &GradeRepo{db}
}

func (gr *GradeRepo) GetSemesterGrades(studentID string, semester int) ([]int, error) {
	stmt := `select Grade from grades where StudentID=? and semester=?`
	rows, err := gr.db.Query(stmt, studentID, semester)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var grades []int
	for rows.Next() {
		var grade int
		if err := rows.Scan(&grade); err != nil {
			return nil, err
		}
		grades = append(grades, grade)
	}

	return grades, err
}

func (gr *GradeRepo) AddGrades(studentID string, subjectID string, Grade int, semester int) error {
	stmt := `insert into grades values(?,?,?,?)`
	_, err := gr.db.Exec(stmt, subjectID, studentID, Grade, semester)

	return err
}
