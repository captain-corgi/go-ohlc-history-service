package v1

import (
	"github.com/labstack/echo/v4"
)

type response struct {
	Error string `json:"error" example:"message"`
}

func errorResponse(c echo.Context, code int, msg string) error {
	return c.JSON(code, response{msg})
}
