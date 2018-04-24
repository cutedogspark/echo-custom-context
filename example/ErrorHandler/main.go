package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cutedogspark/echo-custom-context"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func HTTPErrorHandler(err error, c echo.Context) {
	if he, ok := err.(*ctx.GErrCall); ok {
		err = errors.WithStack(he)
		b, _ := json.Marshal(he.ResponseParams)
		c.JSONBlob(he.HttpStatus, b)
	} else if he, ok := err.(*echo.HTTPError); ok {
		// warp echo error struct
		err = errors.WithStack(he)
		gCtx := ctx.CustomCtx{}
		gErrs := gCtx.GResp().Errors(&ctx.GError{
			Code:    uint(he.Code),
			Message: fmt.Sprintf("%+v", he.Message),
		})
		b, _ := json.Marshal(gErrs.ResponseParams)
		c.JSONBlob(he.Code, b)
	} else {
		// define unknown error message
		err = errors.New("unknown error")
		gCtx := ctx.CustomCtx{}
		gErrs := gCtx.GResp().Errors(&ctx.GError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		b, _ := json.Marshal(gErrs.ResponseParams)
		c.JSONBlob(he.Code, b)
	}
	c.Logger().Error(err)
}

func main() {

	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = HTTPErrorHandler
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(ctx.CustomCtx{c})
		}
	})

	e.GET("/", func(c echo.Context) error {
		// no use error handler function
		return c.(ctx.CustomCtx).GResp(http.StatusOK).Data("Service").Out()
	})

	e.GET("/gerr", func(c echo.Context) error {
		gerrs := ctx.NewGErrors().Append(&ctx.GError{
			Code:         40000001,
			Domain:       "Calendar",
			Reason:       "ResourceNotFoundException",
			Message:      "Resources is not exist",
			LocationType: "database query",
			Location:     "query",
			ExtendedHelp: "http://help-link",
			SendReport:   "http://report.dajui.com/",
		}).Append(&ctx.GError{
			Code:         40000002,
			Domain:       "global",
			Reason:       "required",
			Message:      "Required parameter: part",
			LocationType: "parameter",
			Location:     "part",
		})
		return c.(ctx.CustomCtx).GResp().Errors(*gerrs...).Do()
	})

	e.GET("/echo-error", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, "default echo error handler")
	})

	e.GET("/unknown-error", func(c echo.Context) error {
		return errors.New("Goodbye")
	})

	// Start server
	e.Logger.Fatal(e.Start(":1234"))
}
