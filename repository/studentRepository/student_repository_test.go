package studentsRepository_test

import (
	"regexp"
	studentsRepository "sms/repository/studentRepository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetAllStudentsOfClass(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	repo := studentsRepository.NewStudentRepo(db)

	rows := sqlmock.NewRows([]string{"StudentID", "Name", "RollNumber", "ClassID", "semester"}).
		AddRow("1", "Rohith", "RN1", "C1", 1).
		AddRow("2", "Trisha", "RN2", "C1", 1)

	mock.ExpectQuery(regexp.QuoteMeta("select StudentID,Name,RollNumber,ClassID, semester from students where ClassID=? and semester=?")).
		WithArgs("C1", 1).
		WillReturnRows(rows)

	students, err := repo.GetAllStudentsOfClass("C1", 1)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(students) != 2 {
		t.Errorf("expected 2 students, got %d", len(students))
	}

	if students[0].Name != "Rohith" {
		t.Errorf("expected first student to be Rohith, got %s", students[0].Name)
	}
}

func TestAddStudent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	repo := studentsRepository.NewStudentRepo(db)

	mock.ExpectExec(regexp.QuoteMeta("insert into students values(?,?,?,?,?)")).
		WithArgs("1", "Rohith", "RN1", "C1", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.AddStudent("1", "RN1", "Rohith", "C1", 1)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestGetStudentByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	repo := studentsRepository.NewStudentRepo(db)

	row := sqlmock.NewRows([]string{"StudentID", "Name", "RollNumber", "ClassID", "semester"}).
		AddRow("1", "Rohith", "RN1", "C1", 1)

	mock.ExpectQuery(regexp.QuoteMeta("select StudentID,Name,RollNumber,ClassID,semester from students where StudentID=?")).
		WithArgs("1").
		WillReturnRows(row)

	student, err := repo.GetStudentByID("1")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if student == nil {
		t.Fatal("expected a student, got nil")
	}

	if student.Name != "Rohith" {
		t.Errorf("expected student name Rohith, got %s", student.Name)
	}
}

func TestUpdateStudent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	repo := studentsRepository.NewStudentRepo(db)

	mock.ExpectExec(regexp.QuoteMeta("update students set Name=?,RollNumber=?,ClassID=?,semester=? where StudentID=?")).
		WithArgs("RohithUpdated", "RN1", "C1", 1, "1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateStudent("1", "RohithUpdated", "RN1", "C1", 1)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestGetStudentByRollNumber(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	repo := studentsRepository.NewStudentRepo(db)

	row := sqlmock.NewRows([]string{"StudentID", "Name", "RollNumber", "ClassID", "semester"}).
		AddRow("1", "Rohith", "RN1", "C1", 1)

	mock.ExpectQuery(regexp.QuoteMeta("select StudentID,Name,RollNumber,ClassID,semester from students where RollNumber=?")).
		WithArgs("RN1").
		WillReturnRows(row)

	student, err := repo.GetStudentByRollNumber("RN1")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if student == nil {
		t.Fatal("expected a student, got nil")
	}

	if student.Name != "Rohith" {
		t.Errorf("expected student name Rohith, got %s", student.Name)
	}
}
