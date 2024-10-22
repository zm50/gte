package core

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/zm50/gte/trait"
)

var (
	ErrWorkerPoolStopped = errors.New("worker pool is stopped")
	ErrWorkerPoolFull    = errors.New("worker pool is full")
)

// WorkerPool is a pool of workers that can execute tasks concurrently.
type WorkerPool struct {
	taskChan chan func()
	isStopped atomic.Bool
}

var _ trait.WorkerPool = (*WorkerPool)(nil)

func NewWorkerPool(bufferSize int, workerCount int) *WorkerPool {
	pool := &WorkerPool{
		taskChan: make(chan func(), bufferSize),
	}

	pool.initialize(workerCount)

	return pool
}

// initialize creates workerCount number of workers and starts them.
func (w *WorkerPool) initialize(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go func() {
			for task := range w.taskChan {
				task()
			}
		}()
	}
}

// Push adds a task to the task channel. If the task channel is full, it returns an error.
func (w *WorkerPool) Push(task func()) error {
	if w.isStopped.Load() {
		return ErrWorkerPoolStopped
	}

	select {
		case w.taskChan <- task:
	default:
		return ErrWorkerPoolFull
	}

	return nil
}

// BatchPush adds a batch of tasks to the task channel. If the task channel is full, it returns an error.
func (w *WorkerPool) BatchPush(tasks ...func()) (int, error) {
	for i, task := range tasks {
		return i, w.Push(task)
	}

	return len(tasks), nil
}

// PushWithTimeOut adds a task to the task channel with a timeout. If the task channel is full, it returns an error.
func (w *WorkerPool) PushWithTimeOut(timeout time.Duration, task func()) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case w.taskChan <- task:
	case <-timer.C:
		return ErrWorkerPoolFull
	}

	return nil
}

// BatchPushWithTimeOut adds a batch of tasks to the task channel with a timeout. If the task channel is full, it returns an error.
func (w *WorkerPool) BatchPushWithTimeOut(timeout time.Duration, tasks ...func()) (int, error) {
	for i, task := range tasks {
		err := w.PushWithTimeOut(timeout, task)
		if err != nil {
			return i, err
		}
	}
	
	return len(tasks), nil
}

// PushWithContext adds a task to the task channel with a context. If the task channel is full, it returns an error.
func (w *WorkerPool) PushWithContext(ctx context.Context, task func()) error {
	if w.isStopped.Load() {
		return ErrWorkerPoolStopped
	}

	select {
	case w.taskChan <- task:
	case <-ctx.Done():
		return context.Cause(ctx)
	}

	return nil
}

// BatchPushWithContext adds a batch of tasks to the task channel with a context. If the task channel is full, it returns an error.
func (w *WorkerPool) BatchPushWithContext(ctx context.Context, tasks ...func()) (int, error) {
	for i, task := range tasks {
		err := w.PushWithContext(ctx, task)
		if err != nil {
			return i, err
		}
	}
	
	return len(tasks), nil
}

// Stop stops the worker pool.
func (w *WorkerPool) Stop() {
	if !w.isStopped.CompareAndSwap(false, true) {
		return
	}

	close(w.taskChan)
}
