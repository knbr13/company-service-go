package repositories

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail  = errors.New("duplicate email")
	ErrInvalidPassword = errors.New("invalid password")
)

type UserRepository struct {
	DB *sql.DB
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (u *UserRepository) Insert(user *User) error {
	hashedPass, err := user.hashPass()
	if err != nil {
		return err
	}
	_, err = u.DB.Exec(
		"INSERT INTO users (id, username, email, password_hash) VALUES (?, ?, ?, ?)",
		user.ID,
		user.Username,
		user.Email,
		hashedPass,
	)
	if err != nil {
		if errors.Is(err, &mysql.MySQLError{Number: 1062}) {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (u *User) hashPass() ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
}

func (u *User) Matches(hash []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(u.Password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (u *UserRepository) GetUserByEmail(email string) (*User, error) {
	row := u.DB.QueryRow(
		"SELECT id, username, email, password_hash FROM users WHERE email =?",
		email,
	)
	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
