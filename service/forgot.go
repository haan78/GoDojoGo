package service

import (
	"GoDojoGo/data"

	"github.com/labstack/echo/v5"
)

type UserForgotPasswordServiceType struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
}

func ForgotPasswordService(c *echo.Context) error {
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

func ForgotPasswordSetService(c *echo.Context) error {
	guid := c.Param("guid")
	var msg string = ""
	if guid != "" {
		err := data.SetNewPasswordByGuid(guid)
		if err != nil {
			msg = err.Error()
		}
	} else {
		msg = "No guid where found!"
	}

	if msg == "" {
		c.Render(200, "forgot.html", nil)
	} else {
		data := map[string]any{
			"Message": msg,
		}
		c.Render(200, "error.html", data)
	}
	return nil

}
