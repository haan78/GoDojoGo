package service

import (
	data "GoDojoGo/data"
	lib "GoDojoGo/lib"
	"os"

	"github.com/labstack/echo/v5"
)

func ServiceCreateUser(c *echo.Context) error {
	var req data.UserDetailType
	err := c.Bind(&req)
	if err == nil {
		_, err := data.CreateOrUpdateUser(&req)
		if err == nil {
			err := lib.SendinblueTemplateEmail(os.Getenv("BREVO_API_KEY"), req.Email, 1, map[string]any{
				"UYE_AD": "Ali",
				"URL":    "https://ankarakendo.com",
			})
			return err
		} else {
			return err
		}
	} else {
		return err
	}
}
