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

/*
	activity.GET("/listuser/:activityId", service.ListActivityUserService)
	activity.GET("/adduser/:userId/:activityId", service.AddActivityUserService)
	activity.GET("/deluser/:userId/:activityId", service.DelActivityUserService)
*/

func ListActivityUserService(c *echo.Context) error {
	activityId, err := strconv.Atoi(c.Param("activityId"))
	if err == nil {
		list, err := data.GetUsersOfActivity(int64(activityId))
		if err == nil {
			c.JSON(200, list)
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func PosibleUsersOfActivityService(c *echo.Context) error {
	activityId, err := strconv.Atoi(c.Param("activityId"))
	if err == nil {
		list, err := data.PosibleUsersOfActivity(int64(activityId))
		if err == nil {
			c.JSON(200, list)
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func AddActivityUserService(c *echo.Context) error {
	activityId, err := strconv.Atoi(c.Param("activityId"))
	if err == nil {
		userId, err := strconv.Atoi(c.Param("userId"))
		if err == nil {
			mr, err := data.MonetaryRecordAddByActivity(int64(activityId), int64(userId))
			if err == nil {
				c.JSON(200, mr)
				return nil
			} else {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}
}

func DelActivityUserService(c *echo.Context) error {
	activityId, err := strconv.Atoi(c.Param("activityId"))
	if err == nil {
		userId, err := strconv.Atoi(c.Param("userId"))
		if err == nil {
			err := data.MonetaryRecordDelByActivity(int64(activityId), int64(userId))
			if err == nil {
				c.JSON(200, true)
				return nil
			} else {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}
}
