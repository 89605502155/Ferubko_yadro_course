package repository

import (
	"errors"
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
	testTable := []struct {
		name          string
		user          server.User
		mockBechavior mockBechavior
		expectedError error
	}{
		{
			name: "OK",
			user: server.User{
				Username: "andrey",
				Password: "123",
				Status:   "admin",
			},
			mockBechavior: func(user server.User) {
				query := fmt.Sprintf("INSERT INTO %s ((.+)) VALUES  ((.+),  (.+),  (.+))", usersTable)

				mock.ExpectExec(query).WithArgs(user.Username, user.Password, user.Status).WillReturnResult(
					sqlxmock.NewResult(1, 1))

			},
			expectedError: nil,
		},
		{
			name: "no Name",
			user: server.User{
				Username: "",
				Password: "123",
				Status:   "admin",
			},
			mockBechavior: func(user server.User) {
				query := fmt.Sprintf("INSERT INTO %s ((.+)) VALUES  ((.+),  (.+),  (.+))", usersTable)
				mock.ExpectExec(query).WithArgs(user.Username, user.Password, user.Status).WillReturnError(errors.New("no name"))
			},
			expectedError: errors.New("no name"),
		},
		{
			name: "no Password",
			user: server.User{
				Username: "andrey",
				Password: "",
				Status:   "admin",
			},
			mockBechavior: func(user server.User) {
				query := fmt.Sprintf("INSERT INTO %s ((.+)) VALUES  ((.+),  (.+),  (.+))", usersTable)
				mock.ExpectExec(query).WithArgs(user.Username, user.Password, user.Status).WillReturnError(errors.New("no password"))
			},
			expectedError: errors.New("no password"),
		},
		{
			name: "bad status",
			user: server.User{
				Username: "andrey",
				Password: "123",
				Status:   "admiral",
			},
			mockBechavior: func(user server.User) {
				query := fmt.Sprintf("INSERT INTO %s ((.+)) VALUES  ((.+),  (.+),  (.+))", usersTable)
				mock.ExpectExec(query).WithArgs(user.Username, user.Password, user.Status).WillReturnError(errors.New("bad status"))
			},
			expectedError: errors.New("bad status"),
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			test.mockBechavior(test.user)
			err := r.CreateUser(test.user)
			if errors.Is(test.expectedError, err) && err != nil {
				t.Log(errors.Is(test.expectedError, err), "\n")
				t.Errorf("expected error %v, got %v", test.expectedError, err)
			} else {
				t.Log("Good")
			}
		})
	}

}
