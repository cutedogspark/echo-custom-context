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
			name: "400",
			givenHandler: func(c echo.Context) error {

				var errDate []ctx.ErrorProto

				errCode := 40000001

				errDate = append(errDate, ctx.ErrorProto{
					Domain: "Calendar",
					Reason: "ResourceNotFoundException",
					Message: "Resources is not exist",
					LocationType: "database query",
					Location: "query",
					ExtendedHelp: "http://help-link",
					SendReport: "http://report.dajui.com/",
				})

				errMsg := errDate[0].Message

				errDate = append(errDate, ctx.ErrorProto{
					Domain: "global",
					Reason: "required",
					Message: "Required parameter: part",
					LocationType: "parameter",
					Location: "part",
				})

				return c.(ctx.CustomCtx).Resp(errCode).Error(fmt.Sprintf("%v", errMsg)).Code(errCode).Errors(errDate).Do()
			},
			wantJSON: `{"apiVersion":"1.0","error":{"code":40000001,"message":"Resources is not exist","errors":[{"extended_help":"http://help-link", "send_report":"http://report.dajui.com/", "domain":"Calendar", "reason":"ResourceNotFoundException", "message":"Resources is not exist", "location":"query", "location_type":"database query"},{"message":"Required parameter: part", "location":"part", "location_type":"parameter", "domain":"global", "reason":"required"}]}}`,
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
