package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type User struct {
	UUID  string
	Token string
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

const tokenExp = time.Hour * 3
const secretKey = "supersecretkey"

func CreateUser(uuid string) (*User, error) {
	token, err := buildJWTString(uuid)
	if err != nil {
		return nil, err
	}
	return &User{
		UUID:  uuid,
		Token: token,
	}, nil
}

func GetUser(token string) (*User, error) {
	userID, err := getUserID(token)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &User{UUID: userID, Token: token}, nil
}

func buildJWTString(UUID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: UUID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getUserID(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", err
	}

	return claims.UserID, nil
}

func GenerateUUID() string {
	return uuid.New().String()
}
