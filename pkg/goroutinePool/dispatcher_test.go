package helper

import (
	"context"
	"encoding/json"
	"fmt"
	profile "git.mudu.tv/middleware/go-micro/debug/profile/http"
	log "git.mudu.tv/middleware/go-micro/logger"
	models "git.mudu.tv/youke/concurrent/lib/model"
	"git.mudu.tv/youke/utils/zcontext/rediscontext"
	"testing"
	"time"
)

func TestNewDispatcher(t *testing.T) {

	profile.DefaultAddress = ":6069"
	prof := profile.NewProfile()
	prof.Start()
	dispatcher := NewDispatcher("test", 128, 256)
	audienceId := int64(1000000)
	i := 0
	for true {
		audienceId++
		data := models.HeatMessage{
			Id:        audienceId,
			AccountId: int64(1001),
		}
		if audienceId > 1019999 {
			i++
			audienceId = int64(1000000)
		}
		message, _ := json.Marshal(data)
		job := Job{
			Message: message,
			Handle: func(ctx context.Context, message []byte) error {
				c, ok := rediscontext.Get(ctx, models.RedisKey)
				if ok {
					var str string
					var redisErr error
					var action string
					heartMsg := models.HeatMessage{}
					json.Unmarshal(message, &heartMsg)
					key := fmt.Sprintf("test:%d-%d", heartMsg.AccountId, heartMsg.Id)
					if i == 0 {
						action = "Set"
						ok, redisErr = c.HSet("1001", key, string(message)).Result()
					}
					//} else {
					//	action = "Get"
					//	str, redisErr = c.Get(key).Result()
					//}
					log.Infof("action:%v, key:%s string:%v redisErr:%v", action, key, str, redisErr)
				} else {
					log.Info("get client redisErr")
				}
				time.Sleep(time.Microsecond * 2)
				return nil
			},
		}
		log.Info("send msg", data)
		//按照协程创建规则启动，同一个企业只能在一个协程消费
		dispatcher.Put(job)
	}
	dispatcher.Release()
	t.Log("close :")
}
