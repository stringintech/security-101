package store

import (
	"github.com/stringintech/security-101/model"
)

type UserStore struct {
	users map[string]model.User
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]model.User),
	}
}

func (s *UserStore) Save(user model.User) error {
	s.users[user.Username] = user
	return nil
}

func (s *UserStore) FindByUsername(username string) (model.User, bool) {
	user, exists := s.users[username]
	return user, exists
}
