package service

import (
	"GoDojoGo/data"
	"strconv"

	"github.com/labstack/echo/v5"
)

func GetUserDanService(c *echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err == nil {
		result, err := data.GetUserDan(int64(userID))
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
