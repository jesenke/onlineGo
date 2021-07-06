package http

import "github.com/labstack/echo"

func setAdminRoute(e *echo.Echo)  {
	e.Use(checkSecretMiddleWare)
	e.POST("/jwtToken", getJwtToken)
	e.GET("/onlineList", getJwtToken)
	e.GET("/onlineCount", getJwtToken)
	e.GET("/onlineCount", getJwtToken)
}