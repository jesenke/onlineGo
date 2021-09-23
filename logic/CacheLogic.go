package logic

import (
	"github.com/onlineGo/pkg/memory"
	"github.com/sirupsen/logrus"
)

type CacheLogic struct {
	cache *memory.Memory
}

func NewCacheLogic() *CacheLogic {
	return &CacheLogic{
		memory.NewMemory("rpc"),
	}
}

func (this *CacheLogic) Set() {
	logrus.Println("gogo")
}

func (this *CacheLogic) Get() {
	logrus.Println("gogo")

}

func (this *CacheLogic) Exist() {
	logrus.Println("gogo")

}

func (this *CacheLogic) List() {
	logrus.Println("gogo")

}

func (this *CacheLogic) Count() {
	logrus.Println("gogo")
}

func (this *CacheLogic) Del() {
	logrus.Println("gogo")
}
