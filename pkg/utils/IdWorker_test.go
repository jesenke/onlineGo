package utils

import (
	"testing"
	"time"
)

func TestGetLocalIp(p *testing.T) {
	println(1 << 1)
	//t.Log(GetId())
	//t.Log(GetId())
	//t.Log(GetId())
	//t.Log(GetId())
	//for i := 0; i < 10; i++ {
	//	go func() {
	//		for  {
	//			t.Log(GetId())
	//		}
	//	}()
	//}
	//time.Sleep(30 * time.Second)
	//避免阻塞
	t := time.NewTicker(1 * time.Second)
	row := make(chan interface{})
	go func() {
		<-row
	}()
	//避免阻塞
	select {
	case row <- 1:
		p.Log("success")
	case <-t.C:
		p.Log("fail")
	}
}
