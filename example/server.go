package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"

	"github.com/cutedogspark/echo-custom-context"
)

func main() {
	e := echo.New()

	e.HTTPErrorHandler = ctx.HTTPErrorHandler
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(ctx.CustomCtx{c})
		}
	})

	e.GET("/", func(c echo.Context) error {
		return c.(ctx.CustomCtx).Resp(http.StatusOK).Data("Hello, World!").Do()
	})

	e.GET("v2", func(c echo.Context) error {
		return c.(ctx.CustomCtx).Resp(http.StatusOK).Ver("v2").Data().Do()
	})

	e.GET("v3", func(c echo.Context) error {
		return c.(ctx.CustomCtx).Resp(http.StatusOK).Ver("v3").Data().Do()
	})

	e.GET("/error", func(c echo.Context) error {
		errCode := 40000001
		errMsg := "Error Title"
		errDate := ctx.NewErrors()
		errDate.Add("Error Message 1")
		errDate.Add("Error Message 2")

		return c.(ctx.CustomCtx).Resp(errCode).Error(fmt.Sprintf("%v", errMsg)).Code(errCode).Errors(errDate.Error()).Do()
	})

	e.Logger.Fatal(e.Start(":8080"))
}
