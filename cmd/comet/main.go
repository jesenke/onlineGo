package main

import (
	"github.com/onlineGo/conf"
	"github.com/onlineGo/controller/httpController"
	"github.com/onlineGo/controller/rpcController"
	"github.com/onlineGo/logic"
	"github.com/onlineGo/proto/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

func main() {
	conf.Init()
	//0、服务注册到serverPlat，便于客户端发现
	//1、保存数据[需要断电保护的逻辑:批量数据先直接保存到redis]->mysql记录写入
	//2、定时过期数据[发过期通知到api层]
	//3、数据扩容异步迁移数据发送逻辑[多分片数据一致性(数据设置迁移中的标志，使得数据唯一，迁移完后丢弃，通过rpc服务间自触发)]
	//2种方案：
	//  1、rpc间数据做交换，需要设置标记位【迁出迁入】，避免数据脏读，api读取时直接取新数据（但是数据存在中间态，一定概率的误读），
	//  2、通过api设置rpc服务状态，迁移数据先设置状态为迁出，数据读取认按照迁移前的规则读取迁出、不迁移的2种数据，写按新规则，迁移完后，数据迁移状态刷新掉，
	//当所有分片迁移成功，则通知api,发送迁移完成指令
	//-- 2、迁移数据发出迁移成功后
	//4、数据查询
	//5、数据剔除
	/*难点：tcp长链接，服务注册dns,发现*/
	srv := httpController.RegisterHttp(":8099", conf.GetConfig("ServerKey"), conf.GetConfig("ServerPem"))
	go RpcStart()
	if err := srv.Start(); err != nil {
		logrus.Errorf("the api serve listen err:" + err.Error())
		panic("the api serve start err:" + err.Error())
	}
}

func RpcStart() {
	srv := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{Time: time.Second * 10}))
	reflection.Register(srv)
	listen, err := net.Listen("tcp", ":8091")
	cache := logic.NewCacheLogic()
	handler := rpcController.NewRpcCache(cache)
	proto.RegisterStorageServer(srv, handler)
	if err != nil {
		logrus.Errorf("the rpc serve listen err:" + err.Error())
		panic("the rpc serve listen err:" + err.Error())
	}
	if err := srv.Serve(listen); err != nil {
		logrus.Errorf("the rpc serve start err:" + err.Error())
		panic("the rpc serve start err:" + err.Error())
	}
}
