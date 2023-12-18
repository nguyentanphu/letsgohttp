package models

import (
	"database/sql"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (userModel *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES(?, ?, ?, UTC_TIMESTAMP())`
	_, err = userModel.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (userModel *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hassedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"
	err := userModel.DB.QueryRow(stmt, email).Scan(&id, &hassedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hassedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

func (userModel *UserModel) EmailExists(email string) (bool, error) {
	count := 0
	stmt := `SELECT COUNT(*) FROM users WHERE email=?`
	err := userModel.DB.QueryRow(stmt, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	EmailExists(email string) (bool, error)
	Exists(id int) (bool, error)
	UniqueEmailValidator(fl validator.FieldLevel) bool
}
