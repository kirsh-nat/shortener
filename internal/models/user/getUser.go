package user

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

func GetUser(token string) (*User, error) {
	userID, err := getUserID(token)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &User{UUID: userID, Token: token}, nil
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
