package hash

import (
	"github.com/onlineGo/pkg/utils"
	"hash/crc32"
	"sort"
	"sync"
)

type HashSelector struct {
	base        uint32
	rules       []string
	hash        func(data []byte) uint32
	keys        []uint32
	balancer    map[uint32]string
	newBalancer map[uint32]string
	sync.RWMutex
}

//初始化hash环
func (h *HashSelector) Init(rules []string) {
	h.ReSet(rules)
}

//计算hash节点
func (h *HashSelector) Hash(key string) string {
	h.RLock()
	defer h.RUnlock()
	if h.base == 0 {
		return ""
	}
	hash := h.hash([]byte(key))
	nodeHash := CountNode(h.keys, hash)
	return h.balancer[nodeHash]
}

func (h *HashSelector) ReSet(rules []string) {
	h.Lock()
	defer h.Unlock()
	rules = utils.FilterRepeat(rules)
	h.base = uint32(len(rules))
	h.rules = rules
	h.hash = HashCount
	for _, addr := range rules {
		for i := 0; i < int(h.base); i++ {
			hash := h.hash([]byte(addr))
			h.keys = append(h.keys, hash)
			h.balancer[hash] = addr
		}
	}
}

func (h *HashSelector) CurrentNode() (rules []string) {
	h.RLock()
	defer h.RUnlock()
	return rules
}

func HashCount(key []byte) uint32 {
	return crc32.ChecksumIEEE(key)
}

func CountNode(keys []uint32, hash uint32) uint32 {
	idx := sort.Search(len(keys), func(i int) bool {
		return keys[i] >= hash
	})
	return keys[idx%len(keys)]
}
