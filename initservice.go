package main

import (
	globals "GoDojoGo/deff"
	service "GoDojoGo/service"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type Template struct {
	t *template.Template
}

// Echo v5 Renderer interface
func (r *Template) Render(c *echo.Context, w io.Writer, name string, data any) error {
	return r.t.ExecuteTemplate(w, name, data)
}

func InitService() *echo.Echo {
	e := echo.New()

	e.Static("/static", "static")

	// Initialize and register the renderer
	e.Renderer = &Template{
		t: template.Must(template.ParseGlob("templates/*.html")),
	}

	e.HTTPErrorHandler = service.ServiceErrorHandler

	e.Use(middleware.RequestLogger())

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "GO Dojo Go "+globals.VERSION)
	})

	e.POST("/token", service.CreateTokenReq)

	user := e.Group("/user", service.GetSecLevel(1))

	user.POST("/create", service.ServiceCreateUser)
	user.GET("/all", service.ServiceUserAll)

	user.POST("/photo/:userId", service.SaveUserPhoto)
	user.GET("/photo/:userId", service.GetUserPhoto)

	user.GET("/dan/get/:userId", service.GetUserDanService)
	user.POST("/dan/set", service.SetUserDanService)
	user.POST("/dan/del", service.DelUserDanService)

	user.POST("/password", service.ChangeUserPasswordService)
	user.GET("/depts/:userId", service.MemberDebtsServices)

	forgot := e.Group("/forgot", service.GetSecLevel(0))
	forgot.POST("/", service.ForgotPasswordService)
	forgot.GET("/:guid", service.ForgotPasswordSetService)

	activity := e.Group("/activity", service.GetSecLevel(1))
	activity.GET("/get/:start/:end", service.GetActivityService)
	activity.POST("/set", service.SetActivityService)
	activity.GET("/del/:activityId", service.DelActivityService)
	activity.GET("/listuser/:activityId", service.ListActivityUserService)
	activity.GET("/posibleusers/:activityId", service.PosibleUsersOfActivityService)
	activity.GET("/adduser/:userId/:activityId", service.AddActivityUserService)
	activity.GET("/deluser/:userId/:activityId", service.DelActivityUserService)

	report := e.Group("/report", service.GetSecLevel(1))
	report.GET("/debts", service.DebtsService)
	report.GET("/allmonetaryactions", service.AllMonetaryActionsService)

	monetary := e.Group("/monetary", service.GetSecLevel(1))
	monetary.POST("/expense", service.MonetaryAddExpenseService)
	monetary.GET("/del/:monetaryId", service.MonetaryDelRecordService)
	monetary.POST("/sell", service.MonetarySellService)

	test := e.Group("/test", service.GetSecLevel(0))
	test.GET("/alive", service.TestJustText)
	test.GET("/hello", service.TestHTML)
	return e
}
