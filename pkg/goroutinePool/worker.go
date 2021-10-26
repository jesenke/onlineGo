package helper

import (
	"context"
	log "git.mudu.tv/middleware/go-micro/logger"
)

type StopSignal uint16

const Stop StopSignal = 0

//协程池的最小工作单元，即具体业务处理结构体
type Job struct {
	Message []byte
	Handle  func(ctx context.Context, Message []byte) error
}

type IWorker interface {
	SetCtx(ctx context.Context)
	start()
	put(job Job)
	release()
}

type worker struct {
	ctx        context.Context
	cancel     func()
	workerName string //归属的调度者
	taskPool   chan IWorker
	taskJob    chan Job
	stop       chan StopSignal
}

//资源释放
func (w *worker) release() {
	w.stop <- Stop
}

func (w *worker) SetCtx(ctx context.Context) {
	w.ctx, w.cancel = context.WithCancel(ctx)
	//w.cancel = func() {
	//	//log.Infof("worker release resource: %v", w)
	//	release()
	//	w.cancel()
	//}
}

//投递任务
func (w *worker) put(job Job) {
	w.taskJob <- job
}

//开始任务
func (w *worker) start() {
	go func() {
		var task Job
		for {
			select {
			case task = <-w.taskJob:
				func() {
					log.Infof("%s start:\n ", w.workerName)
					defer func() {
						w.cancel()
						log.Infof("%s close resource:\n", w.workerName)
						if r := recover(); r != nil {
							log.Errorf("worker %s ,execute panic jobData: %v ,err:%v  \n", w.workerName, string(task.Message), r)
						}
					}()
					task.Handle(w.ctx, task.Message)
					return
				}()
				w.taskPool <- w
			case <-w.stop:
				return
			}
		}
	}()
}

func createWorker(name string, workerQueue chan IWorker) IWorker {
	ctx, cancel := context.WithCancel(context.Background())
	w := &worker{
		ctx:        ctx,
		cancel:     cancel,
		workerName: name,
		taskPool:   workerQueue,
		taskJob:    make(chan Job),
		stop:       make(chan StopSignal),
	}
	w.taskPool <- w
	return w
}
