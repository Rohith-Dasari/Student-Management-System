package gradeRepository_test

import (
	"errors"
	"reflect"
	"regexp"
	gradeRepository "sms/repository/gradesRepository"
	"strconv"
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

func TestGetClassAverage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	repo := gradeRepository.NewGradeRepo(db)

	tests := []struct {
		name          string
		classID       string
		semester      int
		mockSetup     func()
		expectedAvg   float64
		expectedError error
	}{
		{
			name:     "Successful retrieval of average",
			classID:  "CS101",
			semester: 1,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`select avg(g.grade) from grades g join students s on s.StudentID=g.StudentID where s.classID=? and g.semester=?`)).
					WithArgs("CS101", 1).
					WillReturnRows(sqlmock.NewRows([]string{"avg_grade"}).AddRow(85.5))
			},
			expectedAvg:   85.5,
			expectedError: nil,
		},
		{
			name:     "Database error during query",
			classID:  "CS103",
			semester: 3,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`select avg(g.grade) from grades g join students s on s.StudentID=g.StudentID where s.classID=? and g.semester=?`)).
					WithArgs("CS103", 3).
					WillReturnError(errors.New("db connection lost"))
			},
			expectedAvg:   0.0,
			expectedError: errors.New("db connection lost"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			avg, err := repo.GetClassAverage(tt.classID, tt.semester)
			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error: %v, got nil", tt.expectedError)
				} else if tt.name == "Scan error" {
					if _, ok := err.(*strconv.NumError); !ok {
						t.Errorf("expected error of type *strconv.NumError, got: %v", err)
					}
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if avg != tt.expectedAvg {
					t.Errorf("expected average: %f, got: %f", tt.expectedAvg, avg)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetToppers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := gradeRepository.NewGradeRepo(db)

	expectedToppers := []gradeRepository.StudentAverage{
		{StudentID: "S001", StudentName: "Alice", Average: 95.5},
		{StudentID: "S002", StudentName: "Bob", Average: 92.0},
		{StudentID: "S003", StudentName: "Charlie", Average: 88.5},
	}

	tests := []struct {
		name            string
		classID         string
		semester        int
		top             int
		mockSetup       func()
		expectedToppers []gradeRepository.StudentAverage
		expectedError   error
	}{
		{
			name:     "Successful retrieval of toppers",
			classID:  "CS101",
			semester: 1,
			top:      3,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"StudentID", "Name", "average"}).
					AddRow("S001", "Alice", 95.5).
					AddRow("S002", "Bob", 92.0).
					AddRow("S003", "Charlie", 88.5)
				mock.ExpectQuery(regexp.QuoteMeta(`select s.StudentID, s.Name, avg(g.grade) as average from grades g
				join students s on s.StudentID=g.StudentID where s.ClassID=? and g.semester=?
				group by s.StudentID, s.Name 
				order by average DESC limit ?`)).
					WithArgs("CS101", 1, 3).
					WillReturnRows(rows)
			},
			expectedToppers: expectedToppers,
			expectedError:   nil,
		},
		{
			name:     "Database error during query",
			classID:  "CS102",
			semester: 2,
			top:      5,
			mockSetup: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`select s.StudentID, s.Name, avg(g.grade) as average from grades g
				join students s on s.StudentID=g.StudentID where s.ClassID=? and g.semester=?
				group by s.StudentID, s.Name 
				order by average DESC limit ?`)).
					WithArgs("CS102", 2, 5).
					WillReturnError(errors.New("db connection lost"))
			},
			expectedToppers: nil,
			expectedError:   errors.New("db connection lost"),
		},
		{
			name:     "No toppers found",
			classID:  "CS103",
			semester: 3,
			top:      10,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"StudentID", "Name", "average"})
				mock.ExpectQuery(regexp.QuoteMeta(`select s.StudentID, s.Name, avg(g.grade) as average from grades g
				join students s on s.StudentID=g.StudentID where s.ClassID=? and g.semester=?
				group by s.StudentID, s.Name 
				order by average DESC limit ?`)).
					WithArgs("CS103", 3, 10).
					WillReturnRows(rows)
			},
			expectedToppers: nil,
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			toppers, err := repo.GetToppers(tt.classID, tt.semester, tt.top)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error: %v, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(toppers, tt.expectedToppers) {
					t.Errorf("expected toppers: %v, got: %v", tt.expectedToppers, toppers)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
