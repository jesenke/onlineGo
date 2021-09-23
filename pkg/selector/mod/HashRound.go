package mod

import (
	"github.com/onlineGo/pkg/utils"
	"sync"
)

type HashRound struct {
	base     uint32
	rules    []string
	balancer map[uint]string
	sync.RWMutex
}

func (h *HashRound) Init(rules []string) {
	h.base = uint32(len(rules))
	h.balancer = make(map[uint]string)
	for k, v := range rules {
		h.balancer[uint(k)] = v
	}
}

func (h *HashRound) Hash(key string) string {
	h.RLock()
	defer h.RUnlock()
	if key == "" {
		return ""
	}
	n := utils.Hash(key, h.base)
	target, ok := h.balancer[n]
	if !ok {
		return ""
	}
	return target
}

func (h *HashRound) ReSet(rules []string) {
	h.Lock()
	defer h.Unlock()
	h.base = uint32(len(rules))
	h.balancer = make(map[uint]string)
	for k, v := range rules {
		h.balancer[uint(k)] = v
	}
}

func (h *HashRound) CurrentNode() (rules []string) {
	h.RLock()
	defer h.RUnlock()

	return h.rules
}
