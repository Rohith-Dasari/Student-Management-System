package gradeRepository_test

import (
	"regexp"
	gradeRepository "sms/repository/gradesRepository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetSemesterGrades(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := gradeRepository.NewGradeRepo(db)
	studentID := "123"
	semester := 1
	rows := sqlmock.NewRows([]string{"Grade"}).
		AddRow(90).
		AddRow(85).
		AddRow(92)

	mock.ExpectQuery(regexp.QuoteMeta(`select Grade from grades where StudentID=? and semester=?`)).
		WithArgs(studentID, semester).
		WillReturnRows(rows)

	grades, err := repo.GetSemesterGrades(studentID, semester)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := []int{90, 85, 92}
	if len(grades) != len(expected) {
		t.Errorf("expected %v grades, got %v", len(expected), len(grades))
	}
	for i := range grades {
		if grades[i] != expected[i] {
			t.Errorf("expected grade %d at index %d, got %d", expected[i], i, grades[i])
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestAddGrades(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := gradeRepository.NewGradeRepo(db)
	studentID := "123"
	subjectID := "sub1"
	grade := 95
	semester := 1

	mock.ExpectExec(regexp.QuoteMeta(`insert into grades values(?,?,?,?)`)).
		WithArgs(subjectID, studentID, grade, semester).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.AddGrades(studentID, subjectID, grade, semester)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
func TestUpdateGrade(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	repo := gradeRepository.NewGradeRepo(db)
	studentID := "123"
	subjectID := "sub1"
	newGrade := 90

	mock.ExpectExec(regexp.QuoteMeta(`update grades set Grade=? where StudentID=? and SubjectID=?`)).
		WithArgs(newGrade, studentID, subjectID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateGrade(studentID, subjectID, newGrade)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
