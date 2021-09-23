package memory

import (
	"sync"
	"time"
)

type Segment struct {
	rows    map[string]*Row
	delRows chan Row
	*sync.RWMutex
}

func (s *Segment) Expire(timeNow time.Time) {
	s.RLock()
	defer s.RUnlock()
	for _, v := range s.rows {
		if timeNow.After(v.expireAt) {
			v.SetStatus(RowStatusDel)
		}
	}
}

func (s *Segment) Set(cell *Row) {
	s.Lock()
	defer s.Unlock()
	val, ok := s.rows[cell.key]
	if ok {
		val.expireAt = cell.expireAt
		val.row = cell
	} else {
		s.rows[cell.key] = cell
	}
}

func (s *Segment) Del(subKey string) {
	s.Lock()
	defer s.Unlock()
	delete(s.rows, subKey)
	return
}

func (s *Segment) DeepCopy() map[string]*Row {
	s.Lock()
	defer s.Unlock()
	cloneTags := make(map[string]*Row)
	for k, v := range s.rows {
		cloneTags[k] = v
	}
	return cloneTags
}

func (s *Segment) CloneValue() map[string]Row {
	s.Lock()
	defer s.Unlock()
	cloneTags := make(map[string]Row)
	for k, v := range s.rows {
		cloneTags[k] = *v
	}
	return cloneTags
}

func (s *Segment) Get(key string) (bool, *Row) {
	s.RLock()
	defer s.RUnlock()
	cell, ok := s.rows[key]
	if !ok {
		return false, nil
	}
	return true, cell
}

//计算需要迁移节点的hash
func (s *Segment) MvOut(keys map[string]bool, channel chan *Row) bool {
	data := s.DeepCopy()
	for _, v := range data {
		if !keys[v.key] {
			continue
		}
		channel <- v
	}
	return true
}

//计算需要迁移节点的hash
func (s *Segment) FlushSeg() {
	data := s.DeepCopy()
	for _, v := range data {
		switch v.status {
		case RowStatusMO:
			//todo
		case RowStatusMI:
		case RowStatusDel:
		}
	}
}
