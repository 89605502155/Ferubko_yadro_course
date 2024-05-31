package service

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"xkcd/pkg/repository"
)

type AuthService struct {
	repo repository.Auth
}
type tokenClaims struct {
	jwt.StandardClaims
	UserId int    `json:"user_id"`
	Status string `json:"status" binding:"required"`
}

func NewAuthService(repo repository.Auth) *AuthService {
	return &AuthService{repo: repo}
}
func (s *AuthService) GenerateToken(username, password string, tokenTTL time.Duration) (string, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		logrus.Println(err, " Generate token first if")
		return "", err
	}
	errPas := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errPas != nil {
		logrus.Println(err, " Generate token second if")
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
		user.Status,
	})
	return token.SignedString([]byte(os.Getenv("JWT_SALT")))
}

func (s *AuthService) ParseToken(accessToken string) (int, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("JWT_SALT")), nil
	})
	if err != nil {
		logrus.Println(err, " Parse token first if")
		return 0, "", err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {

		return 0, "", errors.New("token claims are not of type tokenClaims")
	}
	return claims.UserId, claims.Status, nil
}
