package main

import (
	globals "GoDojoGo/deff"
	service "GoDojoGo/service"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {

	if err := globals.LoadSettings(); err != nil {
		panic(err.Error())
	}

	e := echo.New()

	e.HTTPErrorHandler = service.ServiceErrorHandler

	e.Use(middleware.RequestLogger())

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/token", service.CreateTokenReq)

	user := e.Group("/user", service.GetSecLevel(0))

	user.POST("/create", service.ServiceCreateUser)

	user.POST("/photo/:userId", service.SaveUserPhoto)
	user.GET("/photo/:userId", service.GetUserPhoto)

	user.GET("/dan/get/:userId", service.GetUserDanService)
	user.POST("/dan/set", service.SetUserDanService)
	user.POST("/dan/del", service.DelUserDanService)

	user.GET("/activity/get/:userId", service.GetUserActivitiesService)
	user.POST("/activity/add", service.AddUserActivityService)
	user.POST("/activity/del", service.DelUserActivityService)
	user.POST("/activity/pay", service.UserActivityPaymentService)

	test := e.Group("/test", service.GetSecLevel(0))
	test.POST("/brevo", service.TestBrevoService)

	if err := e.Start(":" + fmt.Sprint(globals.Settings.PORT)); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}

}
