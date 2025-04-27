package user

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
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
