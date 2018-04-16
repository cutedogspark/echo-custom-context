package main

import (
	"github.com/cutedogspark/echo-custom-context"
	"github.com/labstack/echo"
	"net/http"
)

func main() {

	// Echo instance
	e := echo.New()
	e.HideBanner = true
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(ctx.CustomCtx{c})
		}
	})

	e.GET("/", func(c echo.Context) error {
		return c.(ctx.CustomCtx).GResp(http.StatusOK).Data("Service").Do()
	})

	e.GET("/error", func(c echo.Context) error {
		return c.(ctx.CustomCtx).GResp().Errors(&ctx.GError{
			Code:         40001002,
			Reason:       "ParameterInvalid",
			Domain:       "error",
			Message:      "parameter required : id",
			Location:     "id",
			LocationType: "parameter",
		}).Do()
	})

	e.GET("/errors", func(c echo.Context) error {

		ctxErr := ctx.NewGErrors().Append(&ctx.GError{
			Code:         40001003,
			Reason:       "ParameterInvalid",
			Domain:       "validate",
			Message:      "parameter required : id",
			Location:     "id",
			LocationType: "parameter",
		}).Append(&ctx.GError{
			Code:         40001004,
			Reason:       "RecordNotFound",
			Domain:       "repository",
			Message:      "record not found : id",
			Location:     "id",
			LocationType: "user",
		})

		ctxErr.AppendDomain("handler")

		return c.(ctx.CustomCtx).GResp().Errors(*ctxErr...).Do()
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))

}
