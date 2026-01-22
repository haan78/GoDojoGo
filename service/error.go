package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
)

type ServiceErrorType struct {
	Message string
	Status  int
}

func (e *ServiceErrorType) Error() string {
	return e.Message
}

func RaiseServiceError(status int, message string) *ServiceErrorType {
	e := &ServiceErrorType{
		Message: message,
		Status:  status,
	}
	return e
}

func ServiceErrorHandler(c *echo.Context, errOrjinal error) {

	resp, err := echo.UnwrapResponse(c.Response())
	if err == nil {
		if resp.Committed {
			return
		}
	}

	var se *ServiceErrorType

	if errors.As(errOrjinal, &se) {
		_ = c.JSON(se.Status, map[string]string{
			"error": se.Message,
		})
		return
	}

	// fallback
	_ = c.JSON(http.StatusInternalServerError, map[string]string{
		"error": fmt.Sprint(errOrjinal),
	})
}
