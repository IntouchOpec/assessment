package expense

import (
	"net/http"

	"github.com/intouchOpec/assessment/setup_service"
	"github.com/labstack/echo/v4"
)

func GetExpenseDetail(c echo.Context) error {
	var expense Expense
	id := c.Param("id")

	err := setup_service.Db.QueryRow("SELECT id, title, amount, note, tags FROM Expenses WHERE id = $1", id).Scan(&expense.Id, &expense.Title, &expense.Amount, &expense.Note, &expense.Tags)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, expense)
}
