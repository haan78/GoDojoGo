package service

import (
	"GoDojoGo/data"

	"github.com/labstack/echo/v5"
)

func PaymentModelGetAllService(c *echo.Context) error {
	result, err := data.PeymentGetAll()
	if err == nil {
		c.JSON(200, result)
		return nil
	} else {
		return err
	}
}
