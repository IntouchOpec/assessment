package user

import "github.com/golang-jwt/jwt/v4"

type User struct {
	Id           string `json:"id"`
	UserName     string `json:"user_name"`
	PasswordHash string `json:"passwordHash"`
}

type UserBody struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type JwtCustomClaims struct {
	UserName string `json:"username"`
	Id       string `json:"id"`
	jwt.RegisteredClaims
}
