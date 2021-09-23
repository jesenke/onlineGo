package httpController

import (
	"github.com/gin-gonic/gin"
	"github.com/onlineGo/conf"
	"github.com/onlineGo/lib"
	"github.com/onlineGo/logic"
	"github.com/onlineGo/pkg/server"
	"net/http"
)

//调用逻辑
//为了让数据均匀落库，采用按时间hash取模
//任务数控制，根据task数，redis-zset的key不一样
/*客户端调用心跳
redis上线数据结构：uid:status-registerTime:liveId:accountId:桶号【一定规则】, 过期时间（二级缓存住1个3分钟一个5分钟，就可以确保不会应为过期导致误判掉线，确保后面的数据一致）
**服务端http接收心跳,查询redis判断是否上线,是否更新状态：
- 已上线不在同一个组织，根据是否要剔除做不同的逻辑
-- 剔除：更新redis,写入全局队列，2次【1次下线，一次上线】写入全局延时的chan，任务rpc中心执行延时，调用rpc通知下线,
-- 不剔除：新增redis状态,写入全局延时的chan，任务rpc中心执行延时逻辑并通知上线
- 已上线不需要变更状态，直接写入全局延时的chan,任务rpc中心执行延时逻辑
- 已上线需要变更状态, 更新redis, 写入全局的延时chan, 任务rpc中心执行延时并通知变更
- 未上线, 更新redis, 写入全局的延时chan, 任务rpc中心执行延时并通知上线
** task:获取锁后过期定时任务,将过期数据通过rpc调用通知rpc下线接口chan
** rpc通知下线 删除redis,写入全局chan,任务中心执行下线通知（可以不写rpc,在http接口逻辑完成）
- （每s或者100条数据调用）rpc消费延时逻辑：先根据预备的任务数,根据桶号得到的zset的key,将值写入到zset中，score值为过期时间（扩缩容时，根据桶号）
- 扩容缩容接口：根据一定规则将影响到的数据从新分布，此时rpc、task都将数据发往中间态，待数据全部挪移完成，告知task重启重新开始执行
*/
func RegisterHttp(port, certFile, permFile string) *server.Server {
	s := server.NewSimpleServer(port, certFile, permFile)
	h := &AccountController{
		logic: logic.NewApiService(),
	}
	s.AddWrap(func(ctx *gin.Context) *server.Response {
		ticket := ctx.GetHeader("ticket")
		if ticket == "" || ticket != conf.GetConfig("ticket") {
			return server.NewResponse("405", "ticket错误", nil)
		}
		return nil
	})
	s.AddRoute(http.MethodGet, "/api/list", h.List)
	s.AddRoute(http.MethodGet, "/api/exist", h.Exist)
	s.AddRoute(http.MethodGet, "/api/count", h.Count)
	s.AddRoute(http.MethodDelete, "/api/drop", h.Drop)
	s.AddRoute(http.MethodPost, "/api/sign", h.Sign)
	s.AddWrap(func(ctx *gin.Context) *server.Response {
		token := ctx.GetHeader("Token")
		if token == "" {
			return server.NewResponse("405", "token缺失", nil)
		}
		jwt, err := lib.CheckToken(token)
		if err != nil {
			return server.NewResponse("409", "token校验失败", nil)
		}
		ctx.Set("jwt", jwt)
		return nil
	})
	s.AddRoute(http.MethodPost, "/api/beat", h.Beat)

	return s
}

type AccountController struct {
	logic logic.ApiLogic
}

func (this *AccountController) Ticket(ctx *gin.Context) *server.Response {
	return nil
}

func (this *AccountController) Sign(ctx *gin.Context) *server.Response {
	param := make(map[string]string)
	err := ctx.BindJSON(&param)
	if err != nil {
		return server.NewResponse("400", err.Error(), nil)
	}
	res := this.logic.TokenSign(param)
	return &res
}

func (this *AccountController) Beat(ctx *gin.Context) *server.Response {
	val, ok := ctx.Get("jwt")
	if !ok {
		return server.NewResponse("400", "token缺失", nil)
	}
	jwt := val.(lib.JwtClaims)
	res := this.logic.HeartBeat(jwt, ctx.Param("status"))
	return &res
}

func (this *AccountController) List(ctx *gin.Context) *server.Response {
	val, ok := ctx.Get("jwt")
	if !ok {
		return server.NewResponse("400", "token缺失", nil)
	}
	return server.NewResponse("200", "ok", &val)
}

func (this *AccountController) Exist(ctx *gin.Context) *server.Response {
	return nil

}

func (this *AccountController) Count(ctx *gin.Context) *server.Response {
	return nil

}

func (this *AccountController) Drop(ctx *gin.Context) *server.Response {

	return nil
}
