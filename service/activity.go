package service

import (
	"GoDojoGo/data"
	"strconv"

	"github.com/labstack/echo/v5"
)

func SetActivityService(c *echo.Context) error {
	var req data.ActivityBaseType
	err := c.Bind(&req)
	if err == nil {
		id, err := data.SetActivity(&req)
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

func DelActivityService(c *echo.Context) error {
	activityId, err := strconv.Atoi(c.Param("activityId"))
	if err == nil {
		err := data.DelActivity(int64(activityId))
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

func GetActivityService(c *echo.Context) error {
	start := c.Param("start")
	end := c.Param("end")
	list, err := data.GetActivity(start, end)
	if err == nil {
		c.JSON(200, list)
		return nil
	} else {
		return err
	}
}
