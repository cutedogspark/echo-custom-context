package ctx_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"github.com/cutedogspark/echo-custom-context"
)

func TestCustomCtx(t *testing.T) {
	tt := []struct {
		name         string
		givenHandler func(c echo.Context) error
		wantJSON     string
	}{
		{
			name: "200",
			givenHandler: func(c echo.Context) error {
				return c.(ctx.CustomCtx).Resp(http.StatusOK).Data("hello world").Do()
			},
			wantJSON: `{"apiVersion": "1.0", "data": "hello world"}`,
		},
		{
			name: "200 with google json style",
			givenHandler: func(c echo.Context) error {
				return c.(ctx.CustomCtx).GResp(http.StatusOK).Data("hello world").Do()
			},
			wantJSON: `{"apiVersion": "1.0", "data": "hello world"}`,
		},
		{
			name: "400 with google json style",
			givenHandler: func(c echo.Context) error {

				gerrs := ctx.NewGErrors().Append(ctx.GError{
					Code:         40000001,
					Domain:       "Calendar",
					Reason:       "ResourceNotFoundException",
					Message:      "Resources is not exist",
					LocationType: "database query",
					Location:     "query",
					ExtendedHelp: "http://help-link",
					SendReport:   "http://report.dajui.com/",
				}).Append(ctx.GError{
					Code:         40000001,
					Domain:       "global",
					Reason:       "required",
					Message:      "Required parameter: part",
					LocationType: "parameter",
					Location:     "part",
				})

				return c.(ctx.CustomCtx).GResp().Errors(gerrs...).Do()
			},
			wantJSON: `{"apiVersion":"1.0","error":{"code":40000001,"message":"Resources is not exist","errors":[{"extendedHelp":"http://help-link", "sendReport":"http://report.dajui.com/", "domain":"Calendar", "reason":"ResourceNotFoundException", "message":"Resources is not exist", "location":"query", "locationType":"database query"},{"message":"Required parameter: part", "location":"part", "locationType":"parameter", "domain":"global", "reason":"required"}]}}`,
		},
		{
			name: "400 with string errors",
			givenHandler: func(c echo.Context) error {

				errCode := 40000001
				errMsg := "Error Title"
				errDate := ctx.NewErrors()
				errDate.Add("Error Message 1")
				errDate.Add("Error Message 2")

				return c.(ctx.CustomCtx).Resp(errCode).Error(fmt.Sprintf("%v", errMsg)).Code(errCode).Errors(errDate.Error()).Do()
			},
			wantJSON: `{"apiVersion":"1.0","error":{"code":40000001,"message":"Error Title","errors":["Error Message 1","Error Message 2"]}}`,
		},
		{
			name: "400 with custom struct",
			givenHandler: func(c echo.Context) error {
				errs := []interface{}{}
				errs = append(errs, struct {
					Name string `json:"name"`
				}{"peter"})
				return c.(ctx.CustomCtx).Resp(http.StatusOK).Error("this is error message").Errors(errs).Do()
			},
			wantJSON: `{"apiVersion":"1.0","error":{"code":0, "message":"this is error message", "errors":[{"name":"peter"}]}}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(echo.GET, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					return next(ctx.CustomCtx{c})
				}
			}(tc.givenHandler)(c)

			assert.JSONEq(t, tc.wantJSON, rec.Body.String())
		})
	}
}
