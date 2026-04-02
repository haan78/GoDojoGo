package service

import (
	"GoDojoGo/data"

	"github.com/labstack/echo/v5"
)

func MonetaryAddExpenseService(c *echo.Context) error {
	var req data.ExpenseType
	err := c.Bind(&req)
	if err == nil {
		monetaryId, err := data.MonetaryAddExpense(req)
		if err == nil {
			c.JSON(200, monetaryId)
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func MonetarySellService(c *echo.Context) error {
	var req data.SellType
	err := c.Bind(&req)
	if err == nil {
		monetaryId, err := data.Sell(req)
		if err == nil {
			c.JSON(200, monetaryId)
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}
