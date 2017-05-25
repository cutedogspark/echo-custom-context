package main

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
)

var (
	apiVersion = "v2"
)

type CustomCtx struct {
	echo.Context
}

type SuccessResp struct {
	ApiVersion string      `json:"apiVersion"`
	Data       interface{} `json:"data"`
}

type respCall struct {
	c          echo.Context
	httpStatus int
}

func (c CustomCtx) Resp(httpStatus int) *respCall {
	rs := &respCall{
		c:          echo.Context(c),
		httpStatus: httpStatus,
	}
	return rs
}

type dataCall struct {
	c               echo.Context
	httpStatus      int
	responseParams_ SuccessResp
}

func (r *respCall) Data(data interface{}) *dataCall {
	rs := &dataCall{
		c:          r.c,
		httpStatus: r.httpStatus,
		responseParams_: SuccessResp{
			ApiVersion: apiVersion,
			Data:       data,
		},
	}
	return rs
}

// Response Json Format
func (r *dataCall) Do() (err error) {
	b, err := json.Marshal(r.responseParams_)
	if err != nil {
		return err
	}
	return r.c.JSONBlob(r.httpStatus, b)
}

func main() {
	e := echo.New()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(CustomCtx{c})
		}
	})

	e.GET("/test", func(c echo.Context) error {
		return c.(CustomCtx).Resp(http.StatusOK).Data("Hello, World!").Do()
	})

	e.Logger.Fatal(e.Start(":8080"))
}
