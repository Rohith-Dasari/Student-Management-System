package gradeRepository

import "database/sql"

type GradeRepo struct {
	db *sql.DB
}
type StudentAverage struct {
	StudentID   string
	StudentName string
	Average     float64
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

func (gr *GradeRepo) UpdateGrade(studentID string, subjectID string, newGrade int) error {
	stmt := `update grades set Grade=? where StudentID=? and SubjectID=?`
	_, err := gr.db.Exec(stmt, newGrade, studentID, subjectID)
	return err
}

func (gr *GradeRepo) GetAverageGrade(studentID string, semester int) (float64, error) {
	stmt := `select avg(grade) from grades where StudentID=? and semester=?`
	var avg float64
	err := gr.db.QueryRow(stmt, studentID, semester).Scan(&avg)
	return avg, err
}

func (gr *GradeRepo) GetClassAverage(classID string, semester int) (float64, error) {
	stmt := `select avg(g.grade) from grades g join students s on s.StudentID=g.StudentID where s.classID=? and g.semester=?`
	var avg float64
	err := gr.db.QueryRow(stmt, classID, semester).Scan(&avg)
	return avg, err
}

func (gr *GradeRepo) GetToppers(classID string, semester, top int) ([]StudentAverage, error) {

	stmt := `select s.StudentID, s.Name, avg(g.grade) as average from grades g
	join students s on s.StudentID=g.StudentID where s.ClassID=? and g.semester=?
	group by s.StudentID, s.Name 
	order by average DESC limit ?`
	rows, err := gr.db.Query(stmt, classID, semester, top)
	if err != nil {
		return nil, err
	}
	var students []StudentAverage
	for rows.Next() {
		var sa StudentAverage
		if err := rows.Scan(&sa.StudentID, &sa.StudentName, &sa.Average); err != nil {
			return nil, err
		}
		students = append(students, sa)
	}
	return students, nil

}
