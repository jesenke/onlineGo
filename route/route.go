package route

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/onlineGo/conf"
	"github.com/onlineGo/route/http"
)

func Register() *echo.Echo  {
	e := echo.New()
	log := conf.GetLog()
	e.Logger.SetOutput(log.ErrorFile)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: log.AccessFile,
	}))
	http.HttpRegister(e)
	return e
}
