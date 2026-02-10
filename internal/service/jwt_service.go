package service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	CookieName = "auth_token"
)

var (
	InvalidTokenErr = errors.New("invalid token")
	ExpiredTokenErr = errors.New("token has expired")
)

type JWTClaims struct {
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secretKey   []byte
	expiryHours int
}

func NewJWTService() *JWTService {
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("SECRET environment variable is required for JWTService")
	}

	expiryHours := 24 
	if exp := os.Getenv("TOKEN_EXP"); exp != "" {
		var err error
		expiryHours, err = strconv.Atoi(exp)
		if err != nil {
			expiryHours = 24
		}
	}

	return &JWTService{
		secretKey:   []byte(secret),
		expiryHours: expiryHours,
	}
}

func (s *JWTService) GenerateToken(userID int64, name, role string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Name:   name,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "aircraft-system",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ExpiredTokenErr
		}
		return nil, InvalidTokenErr
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, InvalidTokenErr
	}

	return claims, nil
}

func (s *JWTService) GetExpiryDuration() time.Duration {
	return time.Duration(s.expiryHours) * time.Hour
}
