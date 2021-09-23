package utils

import (
	"bytes"
	"strconv"
	"sync"
	"time"
)

type IdWorker struct {
	baseValue             int
	dataCenterId          string //基于ip的随机数
	workerPrefix          string
	sequence              int
	maxNumber 			  int
	minNumber 			  int
	nowNumber 			  int
	buffer 				*bytes.Buffer
	*sync.Mutex
}

var Default = NewIdWorker("dafeaut", 10000,99999)

func GetId() string {
	return Default.NextId()
}
//
func NewIdWorker(workerPrefix string, minNumber, maxNumber int) *IdWorker {
	bufer := bytes.NewBufferString(workerPrefix)
	cet := StringIpToInt(GetLocalIp())
	str := strconv.Itoa(cet)
	bufer.WriteString(":")
	bufer.WriteString(str)
	prefix := bufer.String()
	bufer.Reset()
	worker := IdWorker{
		baseValue: time.Now().Second(),
		dataCenterId: str,
		workerPrefix: prefix,
		sequence: 1,
		maxNumber: maxNumber,
		minNumber: minNumber,
		nowNumber: minNumber,
		Mutex: &sync.Mutex{},
		buffer: bufer,
	}
	return &worker
}

func (i *IdWorker) resetWorker()   {
	i.Lock()
	defer i.Unlock()
	i.nowNumber = i.minNumber
	i.baseValue = time.Now().Second()
}

func (i *IdWorker) toString() string  {
	i.buffer.Reset()
	i.buffer.WriteString(i.workerPrefix)
	i.buffer.WriteString(strconv.Itoa(i.baseValue))

	i.buffer.WriteString(strconv.Itoa(i.nowNumber))
	return i.buffer.String()
}

func (i *IdWorker) NextId() string {
	i.Lock()
	i.nowNumber = i.nowNumber + i.sequence
	if i.nowNumber > i.maxNumber {
		i.Unlock()
		i.resetWorker()
		return i.toString()
	}
	str := i.toString()
	i.Unlock()
	return str
}

