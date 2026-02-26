package service

import (
	data "GoDojoGo/data"
	"fmt"

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

type UserForgotPasswordServiceType struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
}

func UserForgotPasswordService(c *echo.Context) error {
	var req UserForgotPasswordServiceType
	err := c.Bind(&req)
	if err == nil {
		guid, err := data.GenerateUserGuidForForgotPass(req.Email, req.NewPassword)
		if err == nil {
			c.JSON(200, guid)
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func UserForgotPasswordSetService(c *echo.Context) error {
	guid := c.Param("guid")
	if guid != "" {
		err := data.SetNewPasswordByGuid(guid)
		if err == nil {
			c.JSON(200, true)
			return nil
		} else {
			return err
		}
	} else {
		return fmt.Errorf("no guid")
	}
}
