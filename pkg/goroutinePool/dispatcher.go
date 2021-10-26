package helper

import (
	"context"
	"fmt"
	models "git.mudu.tv/youke/concurrent/lib/model"
	"git.mudu.tv/youke/utils/zcontext/mysqlcontext"
	"git.mudu.tv/youke/utils/zcontext/rediscontext"
	"git.mudu.tv/youke/utils/zmysql"
	"git.mudu.tv/youke/utils/zredis"
	rd "github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"sync"
	"time"
)

type IPool interface {
	Put(task Job)
	dispatch()
	getWorkerPool() chan IWorker
	WaitCount(count int)
	TaskDone()
	WaitAll()
	Release()
	TaskQueueLen() int
}

type Dispatcher struct {
	name       string
	workerPool chan IWorker
	jobQueue   chan Job
	stop       chan StopSignal
	MaxIdle    uint64
	MaxRun     uint64
	total      uint64
	wg         sync.WaitGroup
	ctx        context.Context
}

func (p *Dispatcher) dispatch() {
	var task Job
	for {
		select {
		case task = <-p.jobQueue:
			w := <-p.workerPool
			w.put(task)
			p.total++
			if p.total > 20000 {
				p.total = 0
			}
		}
	}
}

func (p *Dispatcher) getWorkerPool() chan IWorker {
	return p.workerPool
}

func (p *Dispatcher) Put(task Job) {
	p.jobQueue <- task
}

func (p *Dispatcher) WaitCount(count int) {
	p.wg.Add(count)
}

func (p *Dispatcher) TaskDone() {
	p.wg.Done()
}

func (p *Dispatcher) WaitAll() {
	p.wg.Wait()
}

func (p *Dispatcher) Release() {
	for i := 0; i < cap(p.workerPool); i++ {
		w := <-p.workerPool
		w.release()
	}
	p.stop <- Stop

}

func (p *Dispatcher) TaskQueueLen() int {
	return len(p.jobQueue)
}

const defaultPoolLen = 256
const runningTaskLen = 256

func NewDispatcher(name string, workPoolLen, taskQueueLen uint64, ctx context.Context) IPool {
	if workPoolLen == 0 {
		workPoolLen = defaultPoolLen
	}
	if taskQueueLen == 0 {
		taskQueueLen = runningTaskLen
	}
	var p IPool = &Dispatcher{
		name:       name,
		workerPool: make(chan IWorker, workPoolLen),
		jobQueue:   make(chan Job, taskQueueLen),
		stop:       make(chan StopSignal),
		wg:         sync.WaitGroup{},
	}
	p.WaitCount(int(workPoolLen))
	for i := 0; i < int(workPoolLen); i++ {
		name := fmt.Sprintf("worker: %s-%d", name, i)
		w := createWorker(name, p.getWorkerPool())
		w.SetCtx(ctx)
		w.start()
		p.TaskDone()
	}
	go p.dispatch()
	return p
}
