package memory

import (
	"net/http"
	"net/rpc"
)

const (
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodGet    = "GET"
	MethodDelete = "DELETE"
)

type MemoryServer struct {
	mem *Memory
	fd  string //写路径
}

func Run(port string) error {
	var s = new(MemoryServer)
	s.mem = new(Memory)
	err := rpc.Register(s)
	if err != nil {
		return err
	}
	rpc.HandleHTTP()
	err = http.ListenAndServe(port, nil)
	if err != nil {
		return err
	}
	return nil
}

//---方法需要传任期来确保是否还是主节点
func (s *MemoryServer) Get(key []byte, result *Result) error {
	return nil
}

func (s *MemoryServer) Del(key []byte, result *Result) error {
	return nil
}

func (s *MemoryServer) List(key []byte, result *Result) error {
	return nil
}

func (s *MemoryServer) Add(param []byte, result *Result) error {
	return nil
}

//依次判断是否存在于内存、磁盘、对应节点（过期数据不会直接删，而是改变状态值到-999）
func (s *MemoryServer) Exist(key, result *Result) error {
	return nil
}

func (s *MemoryServer) Expire(key []byte, result *Result) error {
	return nil
}

//将数据写如磁盘(做副本时使用)
func (s *MemoryServer) Write(data []byte, result *Result) error {
	return nil
}

//将数据从新按照规则去分割副本规则
func (s *MemoryServer) Split(data []byte, result *Result) error {
	return nil
}

//选举处理
func (s *MemoryServer) Select(data []byte, result *Result) error {
	return nil
}
