package http

import (
	"github.com/labstack/echo"
	"github.com/onlineGo/conf"
	"github.com/onlineGo/lib"
	"time"
)

func HttpRegister(e *echo.Echo) {


	e.Group("OnlineApi")
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, time.Now().Format(time.RFC3339Nano))
	})
	//客户端请求
	e.POST("/heartBeat", heartBeat, checkJwtTokenMiddleWare)
	//服务端接口
	setAdminRoute(e)
}



func checkSecretMiddleWare (next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("token") == conf.GetConfig("jwtTokenSign") {
			return c.JSON(401, "no access")
		}
		return next(c)
	}
}

func checkJwtTokenMiddleWare(next echo.HandlerFunc) echo.HandlerFunc  {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(401, "no access")
		}
		jwtClaim, err := lib.CheckToken(token)
		if err != nil  {
			return c.JSON(401, "jwt token not valid")
		}
		if !jwtClaim.VerifyExpiresAt(int64(time.Now().Second()), true) {
			return c.JSON(400, "jwt token expire ")
		}
		c.Set("claim", jwtClaim)
		return next(c)
	}
}

