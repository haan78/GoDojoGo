package service

import (
	"github.com/labstack/echo/v5"
)

type TestBrevoServiceType struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Url   string `json:"url"`
}

func TestJustText(c *echo.Context) error {
	c.JSON(200, "I am alive!")
	return nil
}

func TestHTML(c *echo.Context) error {
	data := map[string]any{
		"Name": "World",
	}
	return c.Render(200, "hello.html", data)
}
