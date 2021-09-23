package wrap

import (
	"github.com/labstack/echo"
	"github.com/onlineGo/conf"
	"github.com/onlineGo/controller/httpController"
	"github.com/onlineGo/lib"
	"github.com/onlineGo/logic"
	"time"
)

type OnlineService struct {
	logic logic.ApiLogic
}

func NewService() OnlineService {
	return OnlineService{
		logic.NewApiService(),
	}
}

func HttpRegister(e *echo.Echo) {
	e.GET("/health", func(c echo.Context) error {
		data := make(map[string]interface{})
		data["code"] = "200"
		data["msg"] = "ok"
		data["data"] = time.Now().Format(time.RFC3339Nano)
		return c.JSON(200, data)
	})
	handle := NewService()
	//客户端请求

	//服务端接口

	return
}

func checkSecretMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("token") == conf.GetConfig("jwtTokenSign") {
			return c.JSON(401, "no access")
		}
		return next(c)
	}
}

func checkJwtTokenMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(401, "no access")
		}
		jwtClaim, err := lib.CheckToken(token)
		if err != nil {
			return c.JSON(401, "jwt token not valid")
		}
		if !jwtClaim.VerifyExpiresAt(int64(time.Now().Second()), true) {
			return c.JSON(400, "jwt token expire ")
		}
		c.Set("claim", jwtClaim)
		return next(c)
	}
}
