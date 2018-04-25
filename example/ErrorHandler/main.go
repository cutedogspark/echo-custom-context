package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cutedogspark/echo-custom-context"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

func HTTPErrorHandler(err error, c echo.Context) {
	if he, ok := err.(*ctx.GErrCall); ok {
		err = errors.WithStack(he)
		b, _ := json.Marshal(he.ResponseParams)
		c.JSONBlob(he.HttpStatus, b)
	} else if he, ok := err.(*ctx.GError); ok {
		err = errors.WithStack(he)
		gErrs := ctx.CustomCtx{}.GResp().Errors(he)
		b, _ := json.Marshal(gErrs.ResponseParams)
		c.JSONBlob(gErrs.HttpStatus, b)
	} else if he, ok := err.(*echo.HTTPError); ok {
		// warp echo error struct
		err = errors.WithStack(he)
		gErrs := ctx.CustomCtx{}.GResp().Errors(&ctx.GError{
			Code:    uint(he.Code),
			Message: fmt.Sprintf("%+v", he.Message),
		})
		b, _ := json.Marshal(gErrs.ResponseParams)
		c.JSONBlob(gErrs.HttpStatus, b)
	} else if _, ok := err.(*validator.InvalidValidationError); !ok {
		var errMsg []string
		for _, err := range err.(validator.ValidationErrors) {
			errMsg = append(errMsg, fmt.Sprintf("%s:%s", err.Field(), err.ActualTag()))
		}
		gErrs := ctx.CustomCtx{}.GResp(http.StatusBadRequest).Errors(&ctx.GError{Code: http.StatusBadRequest, Message: strings.Join(errMsg, ",")})
		b, _ := json.Marshal(gErrs.ResponseParams)
		c.JSONBlob(gErrs.HttpStatus, b)
	} else {
		// define unknown error message
		err = errors.New("unknown error")
		gErrs := ctx.CustomCtx{}.GResp().Errors(&ctx.GError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		b, _ := json.Marshal(gErrs.ResponseParams)
		c.JSONBlob(gErrs.HttpStatus, b)
	}
	c.Logger().Error(err)
}

type (
	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {

	e := echo.New()
	e.HideBanner = true
	e.Validator = &CustomValidator{validator: validator.New()}
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

	e.GET("/gerrs", func(c echo.Context) error {
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

	e.GET("/gerr", func(c echo.Context) error {
		gErr := &ctx.GError{
			Code:         40000001,
			Domain:       "Calendar",
			Reason:       "ResourceNotFoundException",
			Message:      "Resources is not exist",
			LocationType: "database query",
			Location:     "query",
			ExtendedHelp: "http://help-link",
			SendReport:   "http://report.dajui.com/",
		}
		return gErr
	})

	e.GET("/echo-error", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, "default echo error handler")
	})

	e.GET("/unknown-error", func(c echo.Context) error {
		return errors.New("Goodbye")
	})

	e.GET("/validate-error", func(c echo.Context) error {

		type req struct {
			App      string `form:"app"             validate:"required,numeric"`
			Key      string `form:"key"             validate:"required"`
			ClientId int    `form:"clientid"        validate:"required"`
		}
		in := new(req)
		in.App = "God"

		if err := c.Validate(in); err != nil {
			return err
		}
		return c.(ctx.CustomCtx).GResp(http.StatusOK).Data("validate sucess").Out()
	})

	// Start server
	e.Logger.Fatal(e.Start(":1234"))
}
