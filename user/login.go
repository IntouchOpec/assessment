package user

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/intouchOpec/assessment/setup_service"
	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	var user UserBody
	var result User
	err := c.Bind(&user)
	if err != nil {

	}

	if err = c.Validate(user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err = setup_service.Db.QueryRow("SELECT id, username, password_hash FROM users WHERE username = $1", user.Username).Scan(&result.Id, &result.UserName, &result.PasswordHash)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if !CheckPasswordHash(user.Password, result.PasswordHash) {
		return echo.ErrUnauthorized
	}

	claims := &JwtCustomClaims{
		result.UserName,
		result.Id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
