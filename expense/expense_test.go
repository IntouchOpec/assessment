package expense_test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/intouchOpec/assessment/expense"
	"github.com/stretchr/testify/assert"
)

type signInResponse struct {
	Token string `json:"token"`
}

func uri(paths ...string) string {
	host := os.Getenv("BASE_URL")
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func register() string {
	body := bytes.NewBufferString(`{
		"username": "systemTest",
		"password": "password"
	}`)
	req, _ := http.NewRequest("POST", uri("auth/register"), body)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, _ := client.Do(req)
	x, _ := ioutil.ReadAll(res.Body)
	var a signInResponse
	json.Unmarshal(x, &a)
	defer res.Body.Close()
	return a.Token
}

func signInApi() string {
	body := bytes.NewBufferString(`{
		"username": "systemTest",
		"password": "password"
	}`)
	req, _ := http.NewRequest("POST", uri("auth/login"), body)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, _ := client.Do(req)
	x, _ := ioutil.ReadAll(res.Body)
	var a signInResponse
	json.Unmarshal(x, &a)
	if res.StatusCode != http.StatusOK {
		return register()
	}
	defer res.Body.Close()
	return a.Token
}

func request(method string, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", "Bearer "+signInApi())
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

type Response struct {
	*http.Response
	err error
}

func TestCreateExpense(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title": "systemTest",
		"amount": 999,
		"note": "system testing"
		}`)
	var expense expense.Expense

	res := request(http.MethodPost, uri("api/expenses"), body)
	x, err := ioutil.ReadAll(res.Body)
	json.Unmarshal(x, &expense)
	defer res.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, expense.Id)
	assert.Equal(t, "systemTest", expense.Title)
	assert.Equal(t, 999, expense.Amount)
}

func TestExpenseList(t *testing.T) {
	var expenses []expense.Expense
	res := request(http.MethodGet, uri("api/expenses"), nil)
	x, err := ioutil.ReadAll(res.Body)
	json.Unmarshal(x, &expenses)
	defer res.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(expenses), 0)
}

func seedExpense(t *testing.T) expense.Expense {
	var expense expense.Expense
	body := bytes.NewBufferString(`{
		"title": "SystemTestSeed",
		"amount": 100,
		"note": "Seed expense.Expense"
	}`)
	res := request(http.MethodPost, uri("api/expenses"), body)
	x, err := ioutil.ReadAll(res.Body)
	json.Unmarshal(x, &expense)
	defer res.Body.Close()

	if err != nil {
		t.Fatal("can't create uomer:", err)
	}
	return expense
}

func TestExpenseDetail(t *testing.T) {
	var expense expense.Expense
	c := seedExpense(t)
	res := request(http.MethodGet, uri("api/expenses", strconv.Itoa(c.Id)), nil)
	x, err := ioutil.ReadAll(res.Body)
	json.Unmarshal(x, &expense)
	defer res.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.NotEqual(t, 0, expense.Id)
	assert.Equal(t, c.Title, expense.Title)
	assert.Equal(t, c.Amount, expense.Amount)
}

func TestUpdateExpense(t *testing.T) {
	var expense expense.Expense
	c := seedExpense(t)
	body := bytes.NewBufferString(`{
		"title": "SystemTestSeedEdit",
		"amount": 99999,
		"note": "Seed expense.Expense"
	}`)
	res := request(http.MethodPut, uri("api/expenses", strconv.Itoa(c.Id)), body)
	x, err := ioutil.ReadAll(res.Body)
	json.Unmarshal(x, &expense)
	defer res.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.NotEqual(t, 0, expense.Id)
	assert.NotEqual(t, c.Title, expense.Title)
	assert.NotEqual(t, c.Amount, expense.Amount)
}
