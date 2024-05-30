package server

import "errors"

type User struct {
	Id       int    `json:"-" db:"id"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" db:"password_hash"`
	Status   string `json:"status" binding:"required"`
}

func (u User) Validate() error {
	if u.Username == "" || u.Password == "" {
		return errors.New("you give nill username or password ")
	}
	statusArray := [4]string{"", "user", "admin", "content manager"}
	res := true
	for _, v := range statusArray {
		if v == u.Status {
			res = false
			break
		}
	}
	if res {
		return errors.New("unknown status")
	}
	return nil
}
