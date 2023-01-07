package expense

import (
	"github.com/lib/pq"
)

type Expense struct {
	Id     int            `json:"id"`
	Title  string         `name:"title" query:"title" json:"title" validate:"required"`
	Amount int            `name:"amount" json:"amount" validate:"required"`
	Note   string         `json:"note"`
	Tags   pq.StringArray `json:"tags"`
}
