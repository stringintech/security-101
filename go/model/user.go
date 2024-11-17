package model

import "time"

type User struct {
	FullName  string    `json:"full_name"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func (u User) GetUsername() string {
	return u.Username
}

func (u User) GetPassword() string {
	return u.Password
}
