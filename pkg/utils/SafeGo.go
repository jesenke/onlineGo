package utils

import (
	"log"
	"sync"
)

type SafeGo struct {
	buffer chan func()
	num    int
	name   string
	sync.WaitGroup
}

func (s *SafeGo) Put(f func()) {
	s.buffer <- f
}

func (s *SafeGo) SetName(name string) {
	s.name = name
}

func NewSafeGo(num int) (s *SafeGo) {
	buffer := make(chan func())
	o := &SafeGo{
		buffer,
		num,
		"",
		sync.WaitGroup{},
	}
	o.Add(o.num)
	for i := 0; i < o.num; i++ {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Println("safeGo:err", err)
				}
			}()
			o.Done()
			select {
			case f := <-buffer:
				f()
			}
		}()
	}
	return o
}
