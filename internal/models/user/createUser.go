package user

import "github.com/google/uuid"

func CreateUser() (*User, error) {
	uuid := generateUUID()
	token, err := buildJWTString(uuid)
	if err != nil {
		return nil, err
	}
	return &User{
		UUID:  uuid,
		Token: token,
	}, nil
}

func generateUUID() string {
	return uuid.New().String()
}
