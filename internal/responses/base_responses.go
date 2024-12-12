package responses

import (
	"net/http"

	"github.com/andrew-sameh/echo-engine/pkg/errors"
	"github.com/labstack/echo/v4"
)

type Response struct {
	Code    int         `json:"-"`
	Pretty  bool        `json:"-"`
	Data    interface{} `json:"data,omitempty"`
	Message interface{} `json:"message"`
}

// sends a JSON response with status code.
func (a Response) JSON(ctx echo.Context) error {
	if a.Message == "" || a.Message == nil {
		a.Message = http.StatusText(a.Code)
	}

	if err, ok := a.Message.(error); ok {
		if errors.Is(err, errors.DatabaseInternalError) {
			a.Code = http.StatusInternalServerError
		}

		if errors.Is(err, errors.DatabaseRecordNotFound) {
			a.Code = http.StatusNotFound
		}

		a.Message = err.Error()
	}

	if a.Pretty {
		return ctx.JSONPretty(a.Code, a, "\t")
	}

	return ctx.JSON(a.Code, a)
}

// type Error struct {
// 	Code  int    `json:"code"`
// 	Error string `json:"error"`
// }

// type Data struct {
// 	Code    int    `json:"code"`
// 	Message string `json:"message"`
// }

// func Response(c echo.Context, statusCode int, data interface{}) error {
// 	// nolint // context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 	// nolint // context.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
// 	// nolint // context.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization")
// 	return c.JSON(statusCode, data)
// }

// func MessageResponse(c echo.Context, statusCode int, message string) error {
// 	return Response(c, statusCode, Data{
// 		Code:    statusCode,
// 		Message: message,
// 	})
// }

// func ErrorResponse(c echo.Context, statusCode int, message string) error {
// 	return Response(c, statusCode, Error{
// 		Code:  statusCode,
// 		Error: message,
// 	})
// }

// Response in order to unify the returned response structure
