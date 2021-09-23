package memory

import (
	"sync"
	"time"
)

type Row struct {
	id            string //行id
	group         string
	key           string
	row           interface{}
	expireAt      time.Time
	status        int //0:不迁移 1:迁出 2:迁入 -1:删除
	*sync.RWMutex     //避免读写冲突
}

const (
	RowStatusDel = -1
	RowStatusMO  = 1
	RowStatusMI  = 2
	RowStatusOK  = 0
)

func (r *Row) SetStatus(status int) {
	r.RLock()
	defer r.RUnlock()
	r.status = status
}

func (r *Row) SetExpire(expire time.Duration) {
	r.RLock()
	defer r.RUnlock()
	r.expireAt = time.Now().Add(expire)
}

func (r *Row) Del() {
	r.Lock()
	defer r.Unlock()
	r.SetStatus(RowStatusOK)
}
