package expense

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/intouchOpec/assessment/setup_service"
	"github.com/intouchOpec/assessment/user"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func Create(c echo.Context) error {
	var expense Expense

	err := c.Bind(&expense)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.Validate(expense); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	userClaims := c.Get("user").(*jwt.Token)
	claims := userClaims.Claims.(*user.JwtCustomClaims)
	userId := claims.Id

	var expenseId int
	stmt, err := setup_service.Db.Prepare("INSERT INTO expenses (title, amount, note, tags, user_id) values ($1,$2,$3, $4, $5) RETURNING id;")

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = stmt.QueryRow(
		expense.Title, expense.Amount, expense.Note, pq.StringArray(expense.Tags), userId).Scan(&expenseId)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	var result Expense

	err = setup_service.Db.QueryRow("SELECT id, title, amount, note, tags from Expenses where id = $1", expenseId).Scan(&result.Id, &result.Title, &result.Amount, &result.Note, &result.Tags)

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, result)
}
