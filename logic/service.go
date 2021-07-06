package logic

import "github.com/onlineGo/lib"

type LogicHandle struct {

}

type ResponseData struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}


func NewConcurrentService() (s LogicHandle) {
	return s
}

func (l LogicHandle) TokenSign (param map[string]string) (responseData ResponseData) {
	claim := lib.NewJwt(param)
	claim.SetExpireAt(3600)
	token, err := claim.GetToken()
	if err!=nil {
		responseData.Code = 403
		responseData.Msg = err.Error()
		responseData.Data = ""
	}
	responseData.Data = token
	return
}

func (l LogicHandle) HeartBeat (claim lib.JwtClaims, status int) (ResponseData ResponseData ) {
	//todo
	//判断完后用管道暂存数据
	return
}
