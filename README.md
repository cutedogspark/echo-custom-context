# echo-custom-context

Response Json Structure with google JSON Style Guide 

- [Google Json style Guide](https://google.github.io/styleguide/jsoncstyleguide.xml)

### Example Code
```go

func main() {
	e := echo.New()

	e.HTTPErrorHandler = ctx.HTTPErrorHandler
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(ctx.CustomCtx{c})
		}
	})

	e.GET("/", func(c echo.Context) error {
		return c.(ctx.CustomCtx).Resp(http.StatusOK).Data("Hello, World!").Do()
	})

	e.GET("v2", func(c echo.Context) error {
		return c.(ctx.CustomCtx).Resp(http.StatusOK).Ver("v2").Data().Do()
	})

	e.GET("v3", func(c echo.Context) error {
		return c.(ctx.CustomCtx).Resp(http.StatusOK).Ver("v3").Data().Do()
	})

	e.GET("/error", func(c echo.Context) error {
		errCode := 40000001
		errMsg := "Error Title"
		errDate := ctx.NewErrors()
		errDate.Add("Error Message 1")
		errDate.Add("Error Message 2")

		return c.(ctx.CustomCtx).Resp(errCode).Error(fmt.Sprintf("%v", errMsg)).Code(errCode).Errors(errDate.Error()).Do()
	})

	e.Logger.Fatal(e.Start(":8080"))
}

```

### Response

##### V1 Sample

```json
{
    "apiVersion": "v1",
    "data": "Hello, World!"
}
```

##### V2 Sample (empty data)

```json
{
    "apiVersion": "v2",
    "data": []
}
```

##### Error Sample

ref : https://google.github.io/styleguide/jsoncstyleguide.xml?showone=error#error
```json
{
    "apiVersion": "v1",
    "error": {
        "code": 40000001,
        "message": "Err Title",
        "errors": [
            "Error Message 1",
            "Error Message 2"
        ]
    }
}
```

##### Page Not Found
```json
{
    "apiVersion": "v1",
    "error": {
    "code": 404,
        "message": "Not Found"
    }
}
```