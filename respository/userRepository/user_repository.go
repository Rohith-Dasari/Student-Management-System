package userrepository

import (
	"database/sql"
	"sms/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (ur *UserRepo) AddUser(id string, name, email, password string) error {
	stmt := `insert into user values(?,?,?,?,?)`
	_, err := ur.db.Exec(stmt, id, name, email, password, "faculty")
	return err
}

func (ur *UserRepo) GetUserByEmailID(email string) (*models.User, error) {
	stmt := `select UserID, Name, Email, Password, Role from user where Email=?`
	row := ur.db.QueryRow(stmt, email)
	var user models.User
	err := row.Scan(&user.UserID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
