package ctx

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

type grespCall struct {
	ver        string
	c          echo.Context
	httpStatus []int
}

func (c CustomCtx) GResp(httpStatus ...int) *grespCall {
	rs := &grespCall{
		ver:        apiVersion,
		c:          echo.Context(c),
		httpStatus: httpStatus,
	}
	return rs
}

type gdataCall struct {
	c              echo.Context
	httpStatus     int
	responseParams SuccessResp
}

func (r *grespCall) Ver(ver string) *grespCall {
	r.ver = ver
	return r
}

func (r *grespCall) Data(data ...interface{}) *gdataCall {
	var d interface{}
	if len(data) == 0 {
		d = []string{}
	} else {
		d = data[0]
	}

	rs := &gdataCall{
		c:          r.c,
		httpStatus: r.httpStatus[0],
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
func (r *gdataCall) Do(replace ...*strings.Replacer) (err error) {
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

// Google JSON Style error call
type gerrorCall struct {
	c              echo.Context
	httpStatus     int
	responseParams GErrorResponse
}

type gerrorMessage struct {
	Code    uint     `json:"code"`
	Message string   `json:"message"`
	Errors  []*GError `json:"errors,omitempty"`
}

type GErrorResponse struct {
	ApiVersion string        `json:"apiVersion"`
	Error      gerrorMessage `json:"error"`
}

func (r *grespCall) Errors(errs ...*GError) *gerrorCall {
	rs := &gerrorCall{
		c: r.c,
		responseParams: GErrorResponse{
			ApiVersion: apiVersion,
			Error:      gerrorMessage{},
		},
	}

	if len(errs) > 0 {
		if len(r.httpStatus) > 0 {
			rs.httpStatus = r.httpStatus[0]
		} else {
			s, _ := strconv.Atoi(fmt.Sprintf("%d", errs[0].Code)[:3])
			rs.httpStatus = s
		}

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
