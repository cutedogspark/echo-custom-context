package ctx

import (
	"encoding/json"
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

type ErrorProto struct {
	Domain       string `json:"domain,omitempty"`
	Reason       string `json:"reason,omitempty"`
	Message      string `json:"message,omitempty"`
	Location     string `json:"location,omitempty"`
	LocationType string `json:"location_type,omitempty"`
	ExtendedHelp string `json:"extended_help,omitempty"`
	SendReport   string `json:"send_report,omitempty"`
}

type errorMessage struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Errors  []ErrorProto `json:"errors,omitempty"`
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

func (r *errorCall) Errors(errors []ErrorProto) *errorCall {
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
