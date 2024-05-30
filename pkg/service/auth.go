package service

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"xkcd/pkg/repository"
)

type AuthService struct {
	repo repository.Auth
}
type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func NewAuthService(repo repository.Auth) *AuthService {
	return &AuthService{repo: repo}
}
func (s *AuthService) GenerateToken(username, password string, tokenTTL time.Duration) (string, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		return "", err
	}
	errPas := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errPas != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})
	return token.SignedString([]byte(os.Getenv("JWT_SALT")))
}
