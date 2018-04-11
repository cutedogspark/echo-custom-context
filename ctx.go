package ctx

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

var (
	apiVersion = "1.0"
)

type CustomCtx struct {
	echo.Context
}

type SuccessResp struct {
	ApiVersion string      `json:"apiVersion"`
	Data       interface{} `json:"data"`
}

type respCall struct {
	ver        string
	c          echo.Context
	httpStatus int
}

func (c CustomCtx) Resp(httpStatus int) *respCall {
	rs := &respCall{
		ver:        apiVersion,
		c:          echo.Context(c),
		httpStatus: httpStatus,
	}
	return rs
}

type dataCall struct {
	c              echo.Context
	httpStatus     int
	responseParams SuccessResp
}

func (r *respCall) Ver(ver string) *respCall {
	r.ver = ver
	return r
}

func (r *respCall) Data(data ...interface{}) *dataCall {
	var d interface{}
	if len(data) == 0 {
		d = []string{}
	} else {
		d = data[0]
	}

	rs := &dataCall{
		c:          r.c,
		httpStatus: r.httpStatus,
		responseParams: SuccessResp{
			ApiVersion: r.ver,
			Data:       d,
		},
	}
	return rs
}

// Response Json Format
// - replace string when response raw data
// - ex: 	replace := strings.NewReplacer("{PP_KEY}", encryptionKey)
func (r *dataCall) Do(replace ...*strings.Replacer) (err error) {
	b, err := json.Marshal(r.responseParams)
	if err != nil {
		return err
	}
	data := string(b)
	for _, value := range replace {
		data = value.Replace(data)
	}

	return r.c.JSONBlob(r.httpStatus, []byte(data))
}

// error call
type errorCall struct {
	c              echo.Context
	code           int
	httpStatus     int
	responseParams ErrorResponse
}

type errorMessage struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Errors  []interface{} `json:"errors,omitempty"`
}

type ErrorResponse struct {
	ApiVersion string       `json:"apiVersion"`
	Error      errorMessage `json:"error"`
}

func (r *respCall) Error(message string) *errorCall {
	rs := &errorCall{
		c:          r.c,
		httpStatus: r.httpStatus,
		responseParams: ErrorResponse{
			ApiVersion: r.ver,
			Error: errorMessage{
				Message: message,
			},
		}}
	return rs
}

func (r *errorCall) Code(code int) *errorCall {
	r.code = code
	r.responseParams.Error.Code = code
	return r
}

func (r *errorCall) Errors(errors []interface{}) *errorCall {
	r.responseParams.Error.Errors = errors
	return r
}

func (r *errorCall) Do() (err error) {
	b, err := json.Marshal(r.responseParams)
	if err != nil {
		return err
	}
	return r.c.JSONBlob(r.httpStatus, b)
}

// Google JSON Style error call
type gerrorCall struct {
	c              echo.Context
	httpStatus     int
	responseParams GErrorResponse
}

type gerrorMessage struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Errors  []GError `json:"errors,omitempty"`
}

type GErrorResponse struct {
	ApiVersion string        `json:"apiVersion"`
	Error      gerrorMessage `json:"error"`
}

func (c CustomCtx) GError(errs ...GError) *gerrorCall {
	rs := &gerrorCall{
		c: echo.Context(c),
		responseParams: GErrorResponse{
			ApiVersion: apiVersion,
			Error:      gerrorMessage{},
		},
	}

	if len(errs) > 0 {
		s, _ := strconv.Atoi(fmt.Sprintf("%d", errs[0].Code)[:3])
		rs.httpStatus = s
		rs.responseParams.Error.Code = errs[0].Code
		rs.responseParams.Error.Message = errs[0].Message
		rs.responseParams.Error.Errors = errs
	}
	return rs
}

func (r *gerrorCall) Do() (err error) {
	b, err := json.Marshal(r.responseParams)
	if err != nil {
		return err
	}
	return r.c.JSONBlob(r.httpStatus, b)
}
