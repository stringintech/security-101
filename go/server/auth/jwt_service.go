package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stringintech/security-101/model"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	jwtSecret []byte
}

func NewService(jwtSecret []byte) *Service {
	return &Service{
		jwtSecret: jwtSecret,
	}
}

func (s *Service) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

func (s *Service) ComparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *Service) GenerateToken(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24 * 10).Unix(), // 10 days
	})

	return token.SignedString(s.jwtSecret)
}

func (s *Service) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", jwt.ErrInvalidKey
	}

	username, ok := claims["sub"].(string)
	if !ok {
		return "", jwt.ErrInvalidKey
	}

	return username, nil
}
