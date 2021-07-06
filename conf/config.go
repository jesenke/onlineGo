package conf

import (
	"encoding/json"
	rd "github.com/go-redis/redis"
	logrus "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)


type redisDns struct {
	Addr 			string			`json:"addr"`
	Idle			int				`json:"idle"`
	Pwd             string          `json:"password"`
	Active			int				`json:"active"`
	IdleTimeOut		int				`json:"idle_time_out"`
	DB				int				`json:"db"`
}

type LogHandle struct {
	ErrorFile *os.File
	AccessFile *os.File
}

type interConf map[string]string

var redis *rd.Client
var log LogHandle
var interConfig interConf



func Init() {
	initRedis()
	initLog()
	initInter()
}

func GetRedis() *rd.Client {
	return redis
}

func GetLog() LogHandle {
	return log
}

func initLog()  {
	path, err := os.Getwd()
	if err != nil {
		panic("log  init fail")
	}
	logFilePath := path + "/log/"
	log.ErrorFile, err = os.OpenFile(logFilePath+"error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic("log error file create  fail")
	}
	log.AccessFile, err = os.OpenFile(logFilePath+"access.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic("log access file create  fail")
	}
	logrus.SetOutput(log.ErrorFile)
	logrus.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	logrus.SetReportCaller(true)
	return
}

func initInter() {
	main, err := loadConfByte("main")
	if err != nil {
		panic("main conf get fail")
	}
	err = json.Unmarshal(main, &interConfig)
	if err != nil {
		panic("redis conf unmarshal fail")
	}
}

func initRedis()  {
	redisConf, err := loadConfByte("redis")
	if err != nil {
		panic("redis conf get fail")
	}
	dsn := redisDns{}
	err = json.Unmarshal(redisConf, &dsn)
	if err != nil {
		panic("redis conf unmarshal fail")
	}
	option := rd.Options{
		Network: "tcp",
		// host:port address.
		Addr: dsn.Addr,

		// Optional password. Must match the password specified in the
		// requirepass server configuration option.
		Password: dsn.Pwd,
		// Database to be selected after connecting to the server.
		DB: dsn.DB,
		// Maximum number of socket connections.
		// Default is 10 connections per every CPU as reported by runtime.NumCPU.
		PoolSize: dsn.Active,

		// Amount of time after which client closes idle connections.
		// Should be less than server's timeout.
		// Default is 5 minutes. -1 disables idle timeout check.
		IdleTimeout: time.Duration(dsn.IdleTimeOut),
	}
	redis = rd.NewClient(&option)
}

func loadConfByte( fileName string ) (str []byte , err error) {
	path := os.Getenv("CONFIG_PATH")
	if len(path) == 0 {
		path, err = os.Getwd()
		if err != nil {
			return nil, err
		}
		path = path + "/conf"
	}

	str, err = ioutil.ReadFile(path  + "/" + fileName + ".json")
	return
}

func GetConfig(key string) string  {
	value, ok := interConfig[key]
	if ok {
		return value
	}
	return ""
}
