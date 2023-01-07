package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v4"
	"github.com/intouchOpec/assessment/expense"
	"github.com/intouchOpec/assessment/setup_service"
	"github.com/intouchOpec/assessment/user"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"
)

type (
	CustomValidator struct {
		Validator *validator.Validate
	}
)

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	}
	return ""
}

type ApiError struct {
	Field string `json:"field"`
	Msg   string `json:"message"`
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		var ve validator.ValidationErrors
		errors.As(err, &ve)
		out := make([]ApiError, len(ve))
		if errors.As(err, &ve) {
			for i, fe := range ve {
				out[i] = ApiError{fe.Field(), msgForTag(fe.Tag())}
			}
		}
		return echo.NewHTTPError(http.StatusBadRequest, out)
	}
	return nil
}

func NotBlank(fl validator.FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.String:
		return len(strings.TrimSpace(field.String())) > 0
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array:
		return field.Len() > 0
	case reflect.Ptr, reflect.Interface, reflect.Func:
		return !field.IsNil()
	default:
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}

func main() {
	setup_service.InitDB()
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v := validator.New()

	err := v.RegisterValidation("notblank", NotBlank)
	if err != nil {
		log.Fatal("register validation faild.", err)
	}
	e.Validator = &CustomValidator{Validator: v}
	auth := e.Group("/auth")
	auth.POST("/register", user.Register)
	auth.POST("/login", user.Login)
	api := e.Group("/api")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(user.JwtCustomClaims)
		},
		SigningKey: []byte("secret"),
	}
	api.Use(echojwt.WithConfig(config))
	api.GET("/expenses", expense.GetList)
	api.POST("/expenses", expense.Create)
	api.GET("/expenses/:id", expense.GetExpenseDetail)
	api.PUT("/expenses/:id", expense.Update)

	go func() {
		if err := e.Start(os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
