package service

import (
	lib "GoDojoGo/lib"
	"os"

	"github.com/labstack/echo/v5"
)

type TestBrevoServiceType struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Url   string `json:"url"`
}

func TestBrevoService(c *echo.Context) error {
	var req TestBrevoServiceType
	err := c.Bind(&req)
	if err == nil {
		err := lib.SendinblueTemplateEmail(os.Getenv("BREVO_API_KEY"), req.Email, 1, map[string]any{
			"UYE_AD": req.Name,
			"URL":    req.Url,
		})
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
