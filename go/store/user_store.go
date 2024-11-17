package store

import (
	"github.com/stringintech/security-101/model"
	"github.com/stringintech/security-101/server/auth"
)

type UserStore struct {
	users map[string]model.User
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]model.User),
	}
}

func (s *UserStore) Create(user model.User) error {
	s.users[user.Username] = user
	return nil
}

func (s *UserStore) GetUserByUsername(username string) (auth.User, bool) { // safe to cast to model.User
	user, exists := s.users[username]
	return user, exists
}
