package ctx

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func HTTPErrorHandler(err error, c echo.Context) {

	r := CustomCtx{c}
	code := http.StatusInternalServerError
	var msg interface{}
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	}
	if err := r.Resp(code).Error(fmt.Sprintf("%v", msg)).Code(code).Do(); err != nil {
		c.Logger().Error(err)
	}

	/*
		errDate := sri.NewErrors()
		errDate.Add(msg)
		if err := r.Resp(code).Error(fmt.Sprintf("%v", msg)).Code(code).Errors(errDate.Error()).Do(); err != nil {
			c.Logger().Error(err)
		}
			{
				"apiVersion": "v1",
				"error": {
					"code": 404,
					"message": "Not Found",
					"errors": [
						"Not Found"
					]
					}
			}
	*/
}
