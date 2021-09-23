package memory

import (
	"github.com/onlineGo/pkg/selector/hash"
	"github.com/onlineGo/pkg/utils"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

const DataStatus = 0
const MvOutStatus = 1
const MvIntStatus = 2

//有列表需要存储
//避免数据倾斜需要怎么处理
//map 的键存组织(企业组织可以合并),第二个键存数据键,快速判断数据是否存在
type Memory struct {
	id      *utils.IdWorker
	channel chan Row //异步落盘的数据

	MvChannel    chan *Row //状态变化的数据
	data         map[string]*Segment
	name         string
	status       int
	maxGoroutine int
	StopSignal   chan struct{}
	*sync.RWMutex
}

var defaultCell = "default"

var DefaultUsers = NewMemory(defaultCell)

func NewMemory(name string) *Memory {
	data := make(map[string]*Segment)
	cellBuffer := make(chan *Row)
	u := &Memory{
		MvChannel:    cellBuffer,
		data:         data,
		id:           utils.NewIdWorker(name, 10000, 99999),
		name:         name,
		status:       0,
		maxGoroutine: 64,
		StopSignal:   make(chan struct{}),
		RWMutex:      &sync.RWMutex{},
	}
	go u.syncSave()
	go u.loopExpire()
	//启动时加载
	go u.loadFile()
	return u
}

func (m *Memory) loadFile() {

}

func (m *Memory) getMemoryCell(key string) *Segment {
	m.RLock()
	model, ok := m.data[key]
	m.RUnlock()
	if ok {
		return model
	} else {
		saveCellModel := make(map[string]*Row)
		saveCell := &Segment{
			rows:    saveCellModel,
			RWMutex: &sync.RWMutex{},
		}
		m.Lock()
		m.data[key] = saveCell
		m.Unlock()
		return saveCell
	}
}

func (m *Memory) syncSave() {
	for {
		select {
		case cell := <-m.channel:
			//todo	存储落盘
			logrus.WithField("row", cell).Error("write to system")
		}
	}
}

func (m *Memory) Set(group, key string, cell interface{}, expire time.Duration) bool {
	model := m.getMemoryCell(group)
	var row *Row

	ok, row := model.Get(key)
	if ok {
		row.SetExpire(expire)
	} else {
		row = &Row{
			m.id.NextId(),
			group,
			key,
			cell,
			time.Now().Add(expire),
			DataStatus,
			&sync.RWMutex{},
		}
		model.Set(row)
	}
	//避免阻塞
	go func() {
		t := time.NewTicker(1 * time.Second)
		select {
		case m.channel <- *row:
		case <-t.C:
			logrus.WithField("row", *row).Error("Set row timeout")
		}
		return
	}()
	return true
}

func (m *Memory) List(group string) map[string]Row {
	model := m.getMemoryCell(group)
	return model.CloneValue()
}

func (u *Memory) Get(group, key string) (ok bool, data interface{}) {
	model := u.getMemoryCell(group)
	ok, cell := model.Get(key)
	if ok {
		return ok, cell.row
	}
	return false, cell.row
}

func (u *Memory) Del(group, key string) bool {
	u.RLock()
	model, ok := u.data[group]
	if ok {
		u.RUnlock()
		return false
	}
	model.Del(key)
	return true
}

func (u *Memory) loopExpire() {
	t := time.NewTicker(time.Second)
	safeGo := utils.NewSafeGo(u.maxGoroutine)
	for {
		select {
		case <-t.C:
			if u.status > 0 {
				continue
			}
			u.RLock()
			timeNow := time.Now()
			for _, saveModel := range u.data {
				safeGo.Put(func() {
					saveModel.Expire(timeNow)
				})
			}
			u.RUnlock()
		}
	}
}

func (u *Memory) MvInto(cell *Row) {
	saveModel := u.getMemoryCell(cell.group)
	saveModel.Set(cell)
}

func (u *Memory) SetStatus(status int) {
	u.Lock()
	defer u.Unlock()
	u.status = status
}

func (u *Memory) StartMv(addr string, hashKeys []string) bool {
	u.SetStatus(MvOutStatus)
	u.RLock()
	defer u.RUnlock()
	data := make(chan *Row)
	Keys := utils.SliceToMap(hashKeys)
	//todo，怎么确定对应key需要迁移
	hash.HashCount([]byte(addr))

	go func() {
		for {
			select {
			case row := <-data:
				row.SetStatus(RowStatusMO)
				//收集数据异步发送
			}
		}
	}()
	for _, saveModel := range u.data {
		saveModel.MvOut(Keys, data)
	}
	return true
}

func (u *Memory) CloseMv() bool {
	u.SetStatus(DataStatus)
	u.RLock()
	defer u.RUnlock()
	safeGo := utils.NewSafeGo(u.maxGoroutine)
	for _, saveModel := range u.data {
		safeGo.Put(func() {
			saveModel.FlushSeg()
		})
	}
	return true
}
