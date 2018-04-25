package ctx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

func HTTPErrorHandler(err error, c echo.Context) {
	if he, ok := err.(*GErrCall); ok {
		err = errors.WithStack(he)
		b, _ := json.Marshal(he.ResponseParams)
		c.JSONBlob(he.HttpStatus, b)
	} else if he, ok := err.(*GError); ok {
		err = errors.WithStack(he)
		gErrs := CustomCtx{}.GResp().Errors(he)
		b, _ := json.Marshal(gErrs.ResponseParams)
		c.JSONBlob(gErrs.HttpStatus, b)
	} else if he, ok := err.(*echo.HTTPError); ok {
		// warp echo error struct
		err = errors.WithStack(he)
		gErrs := CustomCtx{}.GResp().Errors(&GError{
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
		gErrs := CustomCtx{}.GResp(http.StatusBadRequest).Errors(&GError{Code: http.StatusBadRequest, Message: strings.Join(errMsg, ",")})
		b, _ := json.Marshal(gErrs.ResponseParams)
		c.JSONBlob(gErrs.HttpStatus, b)
	} else {
		// define unknown error message
		err = errors.New("unknown error")
		gErrs := CustomCtx{}.GResp().Errors(&GError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		b, _ := json.Marshal(gErrs.ResponseParams)
		c.JSONBlob(gErrs.HttpStatus, b)
	}
	c.Logger().Error(err)
}
