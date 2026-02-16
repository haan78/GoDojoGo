package service

import (
	data "GoDojoGo/data"

	"github.com/labstack/echo/v5"
)

func ServiceCreateUser(c *echo.Context) error {
	var req data.UserDetailType
	err := c.Bind(&req)
	if err == nil {
		_, err := data.CreateOrUpdateUser(&req)
		return err
	} else {
		return err
	}
}

func ServiceUserAll(c *echo.Context) error {
	list, err := data.UserAll()
	if err == nil {
		c.JSON(200, list)
		return nil
	} else {
		return err
	}
}
