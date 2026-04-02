package service

import (
	"GoDojoGo/data"
	"strconv"

	"github.com/labstack/echo/v5"
)

func CalendarListService(c *echo.Context) error {
	active := c.Param("active")
	begin := c.Param("begin")
	end := c.Param("end")
	list, err := data.CalendarList(begin, end, active)
	if err == nil {
		c.JSON(200, list)
		return nil
	} else {
		return err
	}
}

func CalendarDelService(c *echo.Context) error {
	date := c.Param("date")
	activityId, err := strconv.Atoi(c.Param("activityId"))
	if err == nil {
		err := data.CalendarDel(date, activityId)
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

func CalendarAddRemoveMemeberService(c *echo.Context) error {
	calendarId, err := strconv.Atoi(c.Param("calendarId"))
	if err == nil {
		userId, err := strconv.Atoi(c.Param("userId"))
		if err == nil {
			err := data.CalendarAddRemoveMemeber(int64(calendarId), int64(userId))
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

func CalendarMembersService(c *echo.Context) error {
	calendarId, err := strconv.Atoi(c.Param("calendarId"))
	if err == nil {
		list, err := data.CalendarMembers(int64(calendarId))
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
