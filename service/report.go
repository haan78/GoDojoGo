package service

import (
	"GoDojoGo/data"

	"github.com/labstack/echo/v5"
)

func DebtsService(c *echo.Context) error {
	active := c.Param("active")
	result, err := data.Debts(active)
	if err == nil {
		c.JSON(200, result)
		return nil
	} else {
		return err
	}
}

func AllMonetaryActionsService(c *echo.Context) error {
	start := c.Param("active")
	end := c.Param("active")
	result, err := data.AllMonetaryActions(start, end)
	if err == nil {
		c.JSON(200, result)
		return nil
	} else {
		return err
	}
}
