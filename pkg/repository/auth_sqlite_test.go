package repository

import (
	"fmt"
	"reflect"
	"testing"

	sqlxmock "github.com/zhashkevych/go-sqlxmock"

	server "xkcd"
)

func TestGetUser(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	r := NewUserSQLite(db)

	type mockBechavior func(username string)

	testTable := []struct {
		name          string
		username      string
		mockBechavior mockBechavior
		user          server.User
		expectedError error
		expectedEqual bool
	}{
		{
			name:     "OK",
			username: "andrey",
			mockBechavior: func(username string) {
				rows := sqlxmock.NewRows([]string{"username", "password_hash", "status"}).AddRow("andrey", "123", "admin")
				query := fmt.Sprintf("SELECT (.+) FROM %s WHERE username=(.+)", usersTable)
				mock.ExpectQuery(query).WithArgs(username).WillReturnRows(rows)
			},
			user: server.User{
				Username: "andrey",
				Password: "123",
				Status:   "admin",
			},
			expectedError: nil,
			expectedEqual: true,
		},
		{
			name:     "Bad Password",
			username: "andrey",
			mockBechavior: func(username string) {
				rows := sqlxmock.NewRows([]string{"username", "password_hash", "status"}).AddRow("andrey", "1234", "admin")
				query := fmt.Sprintf("SELECT (.+) FROM %s WHERE username=(.+)", usersTable)
				mock.ExpectQuery(query).WithArgs(username).WillReturnRows(rows)
			},
			user: server.User{
				Username: "andrey",
				Password: "123",
				Status:   "admin",
			},
			expectedError: nil,
			expectedEqual: false,
		},
		{
			name:     "Bad status",
			username: "andrey",
			mockBechavior: func(username string) {
				rows := sqlxmock.NewRows([]string{"username", "password_hash", "status"}).AddRow("andrey", "123", "konnetable")
				query := fmt.Sprintf("SELECT (.+) FROM %s WHERE username=(.+)", usersTable)
				mock.ExpectQuery(query).WithArgs(username).WillReturnRows(rows)
			},
			user: server.User{
				Username: "andrey",
				Password: "123",
				Status:   "admin",
			},
			expectedError: nil,
			expectedEqual: false,
		},
		{
			name:     "Not found",
			username: "taras",
			mockBechavior: func(username string) {
				rows := sqlxmock.NewRows([]string{"username", "password_hash", "status"}).AddRow("andrey", "123", "konnetable")
				query := fmt.Sprintf("SELECT (.+) FROM %s WHERE username=(.+)", usersTable)
				mock.ExpectQuery(query).WithArgs(username).WillReturnRows(rows)
			},
			user: server.User{
				Username: "taras",
				Password: "123e",
				Status:   "admiral",
			},
			expectedError: nil,
			expectedEqual: false,
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			test.mockBechavior(test.username)
			got, err := r.GetUser(test.username)
			if test.expectedError != err {
				t.Errorf("expected error %v, got %v", test.expectedError, err)
			} else {
				if reflect.DeepEqual(test.user, got) != test.expectedEqual {
					t.Errorf("expected user %v, got %v", test.user, got)
				}
			}
		})
	}

}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	r := NewUserSQLite(db)

	type mockBechavior func(user server.User)
}
