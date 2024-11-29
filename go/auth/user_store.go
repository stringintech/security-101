package auth

type User interface {
	GetUsername() string
	GetPassword() string
}

type UserStore interface {
	GetUserByUsername(username string) (User, bool)
}
