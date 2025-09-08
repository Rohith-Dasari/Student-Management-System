package userrepository_test

import (
	"regexp"
	userrepository "sms/repository/userRepository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAddUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	repo := userrepository.NewUserRepo(db)

	mock.ExpectExec(regexp.QuoteMeta("insert into user values(?,?,?,?,?)")).
		WithArgs("1", "Rohith", "rohith@example.com", "hashedpass", "faculty").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.AddUser("1", "Rohith", "rohith@example.com", "hashedpass")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUserByEmailID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	repo := userrepository.NewUserRepo(db)

	rows := sqlmock.NewRows([]string{"UserID", "Name", "Email", "Password", "Role"}).
		AddRow("1", "Rohith", "rohith@example.com", "hashedpass", "faculty")

	mock.ExpectQuery(regexp.QuoteMeta("select UserID, Name, Email, Password, Role from user where Email=?")).
		WithArgs("rohith@example.com").
		WillReturnRows(rows)

	user, err := repo.GetUserByEmailID("rohith@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if user == nil {
		t.Fatal("expected a user, got nil")
	}

	if user.Name != "Rohith" {
		t.Errorf("expected user name Rohith, got %s", user.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
