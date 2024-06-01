package server

import (
	"errors"

	"github.com/sirupsen/logrus"
)

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" db:"password_hash"`
	Status   string `json:"status" binding:"required"`
}

func (u *User) Validate() error {
	if u.Username == "" || u.Password == "" {
		return errors.New("you give nill username or password ")
	}
	statusArray := [3]string{"user", "admin", "content manager"}
	res := true
	for _, v := range statusArray {
		if v == u.Status {
			res = false
			break
		}
	}
	logrus.Println(u.Status)
	if res {
		u.Status = "user"
	}
	return nil
}
