package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type JwtServiceConfig struct {
	Secret             []byte
	ExpirationInterval time.Duration
}

type JwtService struct {
	config JwtServiceConfig
}

func NewJwtService(c JwtServiceConfig) *JwtService {
	return &JwtService{c}
}

func (s *JwtService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

func (s *JwtService) ComparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *JwtService) GenerateToken(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject: user.GetUsername(),
		IssuedAt: &jwt.NumericDate{
			Time: time.Now(),
		},
	})
	return token.SignedString(s.config.Secret)
}

func (s *JwtService) ValidateTokenAndGetUsername(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.config.Secret, nil
	}) //TODO use WithValidMethods!!!
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", jwt.ErrInvalidKey
	}

	issuedAt, err := token.Claims.GetIssuedAt()
	if err != nil {
		return "", jwt.ErrInvalidKey
	}
	if issuedAt == nil {
		return "", jwt.ErrInvalidKey
	}
	if issuedAt.Add(s.config.ExpirationInterval).Before(time.Now()) {
		return "", jwt.ErrTokenExpired
	}

	username, err := token.Claims.GetSubject()
	if err != nil || username == "" {
		return "", jwt.ErrInvalidKey
	}
	return username, nil
}
