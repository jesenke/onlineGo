package http

import (
	"github.com/labstack/echo"
	"github.com/onlineGo/lib"
	"github.com/onlineGo/logic"
)


type OnlineService struct {
	logic logic.LogicHandle
}

func NewController() OnlineService {
	return OnlineService{
		logic.NewConcurrentService(),
	}
}


func (s OnlineService)getJwtToken(ctx echo.Context) error {
	param := make(map[string]string)
	if err := ctx.Bind(&param) ; err != nil {
		return ctx.JSON(200, logic.ResponseData{403,err.Error(), ""})
	}
	returnData := s.logic.TokenSign(param)

	return ctx.JSON(200, returnData)
}

func (s OnlineService)heartBeat(ctx echo.Context) error {
	param := struct {
		AccountId string `json:"account_id"`
		Status int `json:"status"`
	}{}
	if err := ctx.Bind(&param) ; err != nil {
		return ctx.JSON(200, logic.ResponseData{403,err.Error(), ""})
	}
	claims := ctx.Get("claim")
	if claims == nil  {
		return ctx.JSON(200, logic.ResponseData{403,"no correct access", ""})
	}
	claimsJwt := claims.(lib.JwtClaims)
	if claimsJwt.UserAccount != param.AccountId {
		return ctx.JSON(200, logic.ResponseData{403,"account not match access", ""})
	}
	return ctx.JSON(200, s.logic.HeartBeat(claimsJwt, param.Status))
}