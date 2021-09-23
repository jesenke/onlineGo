package rpcController

import (
	"context"
	"github.com/onlineGo/pkg/register"
	"github.com/onlineGo/proto/proto"
)

type RpcRegister struct {
	s *register.SimpleRegistry
}

func NewRpcRegister(s *register.SimpleRegistry) proto.RegisterServer {
	return &RpcRegister{s}
}

func (this *RpcRegister) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	return nil, nil
}
func (this *RpcRegister) GetAddr(ctx context.Context, req *proto.GetAddrRequest) (*proto.GetAddrResponse, error) {
	return nil, nil
}
func (this *RpcRegister) Deregister(ctx context.Context, req *proto.DeregisterRequest) (*proto.DeregisterResponse, error) {
	return nil, nil
}
