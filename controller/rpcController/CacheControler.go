package rpcController

import (
	"context"
	"github.com/onlineGo/logic"
	"github.com/onlineGo/proto/proto"
)

type RpcCache struct {
	s *logic.CacheLogic
}

func NewRpcCache(s *logic.CacheLogic) proto.StorageServer {
	return &RpcCache{s}
}

func (this *RpcCache) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	return nil, nil
}

func (this *RpcCache) Set(ctx context.Context, req *proto.SetRequest) (*proto.SetResponse, error) {
	return nil, nil
}
func (this *RpcCache) Exist(ctx context.Context, req *proto.ExistRequest) (*proto.ExistResponse, error) {
	return nil, nil
}
func (this *RpcCache) Count(ctx context.Context, req *proto.CountRequest) (*proto.CountResponse, error) {
	return nil, nil
}
func (this *RpcCache) List(ctx context.Context, req *proto.CountRequest) (*proto.CountResponse, error) {
	return nil, nil
}
func (this *RpcCache) TouchMv(ctx context.Context, req *proto.TouchMvRequest) (*proto.TouchMvResponse, error) {
	return nil, nil
}
