package register

import (
	"io/ioutil"
	"net"
	"sync"
	"time"
)

// GeeRegistry is a simple pkg center, provide following functions.
// add a server and receive heartbeat to keep it alive.
// returns all alive servers and delete dead servers sync simultaneously.
type SimpleRegistry struct {
	timeout time.Duration
	mu      sync.Mutex // protect following
	servers map[string]map[string]*ServerItem
}

type ServerItem struct {
	Addr  string
	Name  string
	Start time.Time
	Conn *net.TCPConn
}

const (
	defaultPath    = "/online/registry"
	defaultTimeout = time.Minute * 5
)

// New create a registry instance with timeout setting
func New(timeout time.Duration) *SimpleRegistry {
	registry := &SimpleRegistry{
		servers: make(map[string]map[string]*ServerItem),
		timeout: timeout,
	}
	go registry.HeatBeat()
	return registry
}

var DefaultGeeRegister = New(defaultTimeout)

func (r *SimpleRegistry) HeatBeat() {
	//对每个链接做好相关心跳检测
	//避免服务下线了不知道
}

func (s *ServerItem)SendMessage(message string)  {

}

func checkServer(item *ServerItem) bool {
	_, err := item.Conn.Write([]byte("ping"))
	if err != nil {
		return false
	}
	res, err := ioutil.ReadAll(item.Conn)
	if err != nil {
		return false
	}
	if string(res) == "pong" {
		return true
	}
	return false
}

func (r *SimpleRegistry) putServer(v ServerItem) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s := r.servers[v.Name]
	if s == nil {
		r.servers[v.Name] = make(map[string]*ServerItem)
		r.servers[v.Name][v.Addr] = &v
	} else {
		r.servers[v.Name][v.Addr].Start =  time.Now()
	}
	return
}

//找到活着者的服务
func (r *SimpleRegistry) aliveServers() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	var alive []string
	for _, s := range r.servers {
		for _, v := range s {
			if v.Start.Add(r.timeout).After(time.Now()) || checkServer(v) {
				alive = append(alive, v.Addr)
			} else {
				//超时服务不可靠了
				r.removeAddr(v)
			}
		}
	}
	return alive
}

//
func (r *SimpleRegistry) removeAddr(v *ServerItem)  {
	if ok := r.servers[v.Name]; ok == nil {
		return
	}
	//通知其他的服务发现者
	go func() {

	}()
}


