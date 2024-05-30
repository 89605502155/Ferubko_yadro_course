package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	server "xkcd"
)

type UserSQLite struct {
	db *sqlx.DB
}

func NewUserSQLite(db *sqlx.DB) *UserSQLite {
	return &UserSQLite{db: db}
}

func (r *UserSQLite) GetUser(username string) (server.User, error) {
	var user server.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE username=$1", usersTable)
	err := r.db.Get(&user, query, username)
	return user, err
}

func (r *UserSQLite) CreateUser(user server.User) error {
	query := fmt.Sprintf("INSERT INTO %s (username, password_hash,status) VALUES  ($1,  $2,  $3)", usersTable)
	_, err := r.db.Exec(query, user.Username, user.Password, user.Status)
	return err
}
