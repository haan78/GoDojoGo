package service

import (
	data "GoDojoGo/data"

	"github.com/labstack/echo/v5"
)

func ServiceCreateUser(c *echo.Context) error {

	type ServiceCreateUserResponse struct {
		UserId int64 `json:"user_id"`
	}

	var req data.UserDetailType
	err := c.Bind(&req)
	if err == nil {
		uid, err := data.CreateOrUpdateUser(&req)
		if err == nil {
			var res ServiceCreateUserResponse
			res.UserId = uid
			c.JSON(200, res)
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

type ChangeUserPasswordServiceType struct {
	Email string `json:"email"`
	Old   string `json:"old"`
	New   string `json:"new"`
}

func ChangeUserPasswordService(c *echo.Context) error {
	var req ChangeUserPasswordServiceType
	err := c.Bind(&req)
	if err == nil {
		return data.ChangeUserPassword(req.Email, req.Old, req.New)
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
