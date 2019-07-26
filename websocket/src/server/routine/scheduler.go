package routine

import (
	"context"
	"fmt"
	"lib/logger"
	"lib/uuid"
	"runtime/debug"
	"time"
)

var routinePool *Routine

func Init(maxGoRoutineCount int) error {
	if maxGoRoutineCount < 1 {
		maxGoRoutineCount = 1
	}
	routinePool = &Routine{
		task:     make(chan *Task, maxGoRoutineCount),
		sem:      make(chan struct{}, maxGoRoutineCount),
		maxCount: int64(maxGoRoutineCount),
	}
	go routinePool.Schedule()
	return nil
}

func Start(ctx context.Context, f TaskFunc) error {
	if routinePool == nil {
		return fmt.Errorf("routine not init")
	}
	t := Task{
		ctx:      ctx,
		taskId:   uuid.MakeUUId(),
		taskFunc: f,
	}
	return routinePool.AddTask(&t)
}

type TaskFunc func(t *Task) error

type Task struct {
	taskId   string
	ctx      context.Context
	taskFunc TaskFunc
}

func (t *Task) GetTaskId() string {
	return t.taskId
}

func (t *Task) execute() (err error) {
	defer func(begin time.Time) {
		if pErr := recover(); pErr != nil {
			logger.Error("[Task.Execute] #PANIC# task execute panic. taskId:%s, panic:%v, stack:%s", t.taskId, err, string(debug.Stack()))
			err = fmt.Errorf("panic err %v", pErr)
		}
		if err != nil {
			logger.Error("[Task.Execute] task execute error. taskId:%s, err:%s", t.taskId, err.Error())
		}
		logger.Info("[Task.Execute] task execute over. taskId:%s, timeUsed:%v", t.taskId, time.Since(begin))
	}(time.Now())
	select {
	case <-t.ctx.Done():
		return fmt.Errorf("task execute timeout")
	default:
	}
	return t.taskFunc(t)
}

type Routine struct {
	task chan *Task
	sem  chan struct{}
	//monitor data
	curRunCount    int64
	totalRunCount  int64
	totalFailCount int64
	maxCount       int64
}

func (r *Routine) Schedule() {
	for {
		select {
		case t := <-r.task:
			go func() {
				defer func() {
					r.sem <- struct{}{}
				}()
				err := t.execute()
				if err != nil {
					r.totalFailCount++
					logger.Error("[Schedule] task exec error. taskId:%s, err:%s", t.taskId, err.Error())
				}
				r.curRunCount--
				r.totalRunCount++
			}()
		}
	}
}

func (r *Routine) AddTask(t *Task) error {
	select {
	case r.sem <- struct{}{}:
		r.curRunCount++
		r.task <- t
	case <-t.ctx.Done():
		return fmt.Errorf("routine add task timeout")
	}
	return nil
}
