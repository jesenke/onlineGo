package register

import (
	"bytes"
	"errors"
	"math/rand"
	"sync"
	"time"
)

const salt = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

//app激活表
//app注册信息
//store需要系统级的秘钥,系统启动时自动生成或者配置文件加载,放在同一wrap
type AppRecord struct {
	AppId  int    //自然数10000开始
	Domain string //应用域:只有在这个域名下使用才能通过网关安全校验
	secret string //一段自然生成的秘钥,内部根据过期时间生成临时密钥使用，可通过秘钥破解破解
}

//app注册表
type AppTable struct {
	NextNumber int
	Table      map[string]AppRecord
	password   string
	sync.RWMutex
}

var appRegister *AppTable

func init() {
	appRegister = newAppTable()
	//也需要定时落盘或者同步磁盘数据
}

func newAppTable() *AppTable {
	return &AppTable{
		NextNumber: 10000,
		Table:      make(map[string]AppRecord),
	}
}

func (app *AppTable) incr() int {
	app.NextNumber = app.NextNumber + 1
	return app.NextNumber
}

func (app *AppTable) salt() string {
	rand.NewSource(time.Now().UnixNano()) // 产生随机种子
	var s bytes.Buffer
	saltLen := int64(len(salt))
	for i := 0; i < 6; i++ {
		s.WriteByte(salt[rand.Int63()%saltLen])
	}
	return s.String()
}

func (app *AppTable) add(domain string) *AppRecord {
	_, ok := app.Table[domain]
	if ok {
		//存在就返回0说添加不成功
		return nil
	}
	app.Lock()
	defer app.Unlock()
	record := AppRecord{
		app.incr(),
		domain,
		app.salt(),
	}
	app.Table[domain] = record
	return &record
}

func Add(domain string) (*AppRecord, error) {
	record := appRegister.add(domain)
	if record == nil {
		return nil, errors.New("domain " + domain + " have exist")
	}
	return record, nil
}
