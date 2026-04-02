package service

import (
	"GoDojoGo/data"

	"github.com/labstack/echo/v5"
)

func ActivityInsertService(c *echo.Context) error {
	var req data.ActivityBaseType
	err := c.Bind(&req)
	if err == nil {
		id, err := data.ActivityInsert(&req)
		if err == nil {
			c.JSON(200, id)
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func ActivityUpdateService(c *echo.Context) error {
	var req data.ActivityBaseType
	err := c.Bind(&req)
	if err == nil {
		err := data.ActivityUpdate(&req)
		if err == nil {
			c.JSON(200, true)
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func ActivityListService(c *echo.Context) error {
	active := c.Param("active")
	list, err := data.ActivityList(active)
	if err == nil {
		c.JSON(200, list)
		return nil
	} else {
		return err
	}
}
