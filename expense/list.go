package expense

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/intouchOpec/assessment/setup_service"
	"github.com/intouchOpec/assessment/user"
	"github.com/labstack/echo/v4"
)

func GetList(c echo.Context) error {
	result := make([]*Expense, 0)
	userClaims := c.Get("user").(*jwt.Token)
	claims := userClaims.Claims.(*user.JwtCustomClaims)
	userId := claims.Id

	rows, err := setup_service.Db.Query("SELECT id, title, amount, note, tags FROM Expenses WHERE user_id = $1", userId)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		expense := new(Expense)

		if err := rows.Scan(&expense.Id, &expense.Title, &expense.Amount, &expense.Note, &expense.Tags); err != nil {
			return err
		}
		result = append(result, expense)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}
