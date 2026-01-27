package service

import (
	"GoDojoGo/data"
	"strconv"

	"github.com/labstack/echo/v5"
)

func AddUserActivityService(c *echo.Context) error {
	var req data.AddUserActivityType
	err := c.Bind(&req)
	if err == nil {
		id, err := data.AddUserActivity(&req)
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

func DelUserActivityService(c *echo.Context) error {
	useractivity_id, err := strconv.Atoi(c.Param("useractivity_id"))
	if err == nil {
		err := data.DelUserActivity(int64(useractivity_id))
		return err
	} else {
		return err
	}
}

func GetUserActivitiesService(c *echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err == nil {
		result, err := data.GetUserActivities(int64(userID))
		if err == nil {
			c.JSON(200, result)
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func UserActivityPaymentService(c *echo.Context) error {
	var req data.UserActivityPaymentType
	err := c.Bind(&req)
	if err == nil {
		err := data.UserActivityPayment(&req)
		return err
	} else {
		return err
	}
}
