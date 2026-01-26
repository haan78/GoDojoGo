package main

import (
	service "GoDojoGo/service"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {

	e := echo.New()

	e.HTTPErrorHandler = service.ServiceErrorHandler

	e.Use(middleware.RequestLogger())

	if err := godotenv.Load(); err != nil {
		e.Logger.Error("failed to start server", "error", err)
		return
	}

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/token", service.CreateTokenReq)

	user := e.Group("/user", service.TokenAuthMiddleware)

	user.POST("/create", service.ServiceCreateUser)
	user.POST("/photo/:userId", service.SaveUserPhoto)
	user.GET("/photo/:userId", service.GetUserPhoto)
	user.GET("/dan/:userId", service.GetUserDanService)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}

}
