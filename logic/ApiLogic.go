package logic

import (
	"github.com/onlineGo/lib"
	"github.com/onlineGo/pkg/selector/hash"
	"github.com/onlineGo/pkg/server"
	"github.com/sirupsen/logrus"
)

type ApiLogic struct {
	selector *hash.HashSelector
}

func NewApiService() (s ApiLogic) {
	return s
}

func (l *ApiLogic) TokenSign(param map[string]string) (responseData server.Response) {
	//同一个人重复认证签名时，创建时间不变
	claim := lib.NewJwt(param)
	claim.SetExpireAt(3600)
	token, err := claim.GetToken()
	if err != nil {
		responseData.Code = "403"
		responseData.Msg = err.Error()
	}
	responseData.Data = token
	return
}

func (l *ApiLogic) HeartBeat(claim lib.JwtClaims, status string) (ResponseData server.Response) {
	//todo
	//判断完后用管道暂存数据
	logrus.Println("gogo")
	return
}

func (l *ApiLogic) OnlineCount(param map[string]string) (ResponseData server.Response) {
	//todo
	//判断完后用管道暂存数据
	logrus.Println("gogo")
	return
}

func (l *ApiLogic) OnlineList(param map[string]string) (ResponseData server.Response) {
	//todo
	//判断完后用管道暂存数据
	logrus.Println("gogo")
	return
}

//扩容时迁移数据接口[数据存储逻辑收到扩容信号后，将所有数据从新请求到api主机，api主机根据新规则从新分配]
func (l *ApiLogic) MvOnlineList(param map[string]string) (ResponseData server.Response) {
	logrus.Println("gogo")
	return
}

func (l *ApiLogic) RegisterNode(nodeAddr string) (ResponseData server.Response) {
	addrs := append(l.selector.CurrentNode(), nodeAddr)
	l.selector.ReSet(addrs)
	return
}
