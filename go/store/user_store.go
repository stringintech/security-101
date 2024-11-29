package store

import (
	"fmt"
	"github.com/stringintech/security-101/auth"
	"github.com/stringintech/security-101/model"
	"sync"
)

type UserStore struct {
	users sync.Map
}

func NewUserStore() *UserStore {
	return &UserStore{}
}

func (s *UserStore) Create(user model.User) error {
	if _, loaded := s.users.LoadOrStore(user.Username, user); loaded {
		return fmt.Errorf("username already exists") //TODO introduce error types
	}
	return nil
}

func (s *UserStore) GetUserByUsername(username string) (auth.User, bool) {
	if value, ok := s.users.Load(username); ok {
		return value.(model.User), true
	}
	return model.User{}, false
}
