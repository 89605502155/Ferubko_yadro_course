package service

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	server "xkcd"
	"xkcd/pkg/repository"
)

type AuthService struct {
	repo repository.Auth
}
type tokenClaims struct {
	jwt.StandardClaims
	Username string `json:"username" binding:"required"`
	Status   string `json:"status" binding:"required"`
	Type     int    `json:"Type"`
}

func NewAuthService(repo repository.Auth) *AuthService {
	return &AuthService{repo: repo}
}

func parseOneToken(tokenTime time.Duration, userName, userStatus, salt string, tokenType int) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTime).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userName,
		userStatus,
		tokenType,
	})
	access, err := accessToken.SignedString([]byte(salt))
	if err != nil {
		logrus.Println(err, " Generate token second access")
		return "", err
	}
	return access, nil
}
func (s *AuthService) GenerateToken(userInput server.UserEntity, accessTime time.Duration, refreshTime time.Duration) (string, string, error) {
	user, err := s.repo.GetUser(userInput.Username)
	if err != nil {
		logrus.Println(err, " Generate token first if")
		return "", "", err
	}
	errPas := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
	if errPas != nil {
		logrus.Println(err, " Generate token second if")
		return "", "", err
	}
	access, err := parseOneToken(accessTime, user.Username, user.Status, os.Getenv("JWT_SALT_ACCESS"), 1)
	if err != nil {
		logrus.Println(err, " Generate token second access")
		return "", "", err
	}

	refresh, err := parseOneToken(refreshTime, user.Username, user.Status, os.Getenv("JWT_SALT_REFRESH"), 2)
	if err != nil {
		logrus.Println(err, " Generate token second refresh")
		return "", "", err
	}
	return access, refresh, nil
}

func parseOneClaims(str, salt string) (*jwt.Token, error) {
	token_, err := jwt.ParseWithClaims(str, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(salt), nil
	})
	return token_, err
}

func (s *AuthService) ParseToken(str string) (string, string, error) {
	accessToken, err := parseOneClaims(str, os.Getenv("JWT_SALT_ACCESS"))
	if err != nil {
		logrus.Println(err, " Parse token first if")
		return "", "", err
	}
	claims, ok := accessToken.Claims.(*tokenClaims)
	if !ok {
		return "", "", errors.New("token claims are not of type tokenClaims")
	}
	if claims.Type != 1 {
		logrus.Println("claims type", claims.Type)
		return "", "", errors.New("token claims are not of access type")
	}
	return claims.Username, claims.Status, nil
}

func (s *AuthService) ParseRefreshToken(str string, accessTime time.Duration) (string, string, string, error) {
	refreshToken, err := parseOneClaims(str, os.Getenv("JWT_SALT_REFRESH"))
	if err != nil {
		logrus.Println(err, " Parse token first if")
		return "", "", "", err
	}
	claims, ok := refreshToken.Claims.(*tokenClaims)
	if !ok {
		return "", "", "", errors.New("token claims are not of type tokenClaims")
	}
	access, err := parseOneToken(accessTime, claims.Username, claims.Status, os.Getenv("JWT_SALT_ACCESS"), 1)
	if err != nil {
		logrus.Println(err, " Generate token second access")
		return "", "", "", err
	}
	return claims.Username, claims.Status, access, nil
}

func (s *AuthService) CreateUser(user server.User) error {
	logrus.Println("start reg")
	pass := user.Password
	coast, _ := strconv.Atoi(os.Getenv("BCRYPT_COAST"))
	logrus.Println("start coast ", coast)
	newPassword, _ := bcrypt.GenerateFromPassword([]byte(pass), coast)
	logrus.Println(user.Password, "  ", newPassword)
	user.Password = string(newPassword)
	logrus.Println(user.Password)
	return s.repo.CreateUser(user)
}
