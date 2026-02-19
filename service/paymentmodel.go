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

func PaymentModelSetAllService(c *echo.Context) error {
	var req []data.PatmentModelType
	err := c.Bind(&req)
	if err == nil {
		err := data.PaymentSetAll(req)
		return err
	} else {
		return err
	}
}
