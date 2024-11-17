package server

import "github.com/stringintech/security-101/model"

type UserStore interface {
	GetUserByUsername(username string) (model.User, bool)
}
