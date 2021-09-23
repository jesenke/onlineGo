package main

import (
	"github.com/onlineGo/conf"
	"github.com/onlineGo/controller/rpcController"
	"github.com/onlineGo/pkg/register"
	"github.com/onlineGo/proto/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

func main() {
	conf.Init()
	srv := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{Time: time.Second * 10}))
	reflection.Register(srv)
	listen, err := net.Listen("tcp", ":8088")
	regsiter := register.New(10 * time.Second)
	handler := rpcController.NewRpcRegister(regsiter)
	proto.RegisterRegisterServer(srv, handler)
	if err != nil {
		panic("the rpc serve start err:" + err.Error())
	}
	srv.Serve(listen)
}
