package expense

import (
	"net/http"

	"github.com/intouchOpec/assessment/setup_service"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func Update(c echo.Context) error {
	var result Expense
	var expense Expense
	err := c.Bind(&expense)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if err = c.Validate(expense); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	id := c.Param("id")
	sqlStatement := `
		UPDATE expenses
		SET title = $2, amount = $3, note = $4, tags = $5
		WHERE id = $1 RETURNING id, title, amount, note, tags;`
	err = setup_service.Db.QueryRow(sqlStatement,
		id, expense.Title, expense.Amount,
		expense.Note, pq.StringArray(expense.Tags)).Scan(
		&result.Id, &result.Title, &result.Amount, &result.Note, &result.Tags)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}
