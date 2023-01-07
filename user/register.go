package user

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/intouchOpec/assessment/setup_service"
	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	var user UserBody
	var userId string
	err := c.Bind(&user)
	if err != nil {

	}
	if err = c.Validate(user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	passwordHash, err := HashPassword(user.Password)

	stmt, err := setup_service.Db.Prepare("INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id")
	err = stmt.QueryRow(
		user.Username, passwordHash).Scan(&userId)
	claims := &JwtCustomClaims{
		user.Username,
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"token": t,
	})
}
