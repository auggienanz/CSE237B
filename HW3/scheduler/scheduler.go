package scheduler

import (
	"log"
	"CSE237B/HW3/task"
	"CSE237B/HW3/worker"
)

// Scheduler dispatches tasks to workers
type Scheduler struct {
	TaskChan      chan *task.Task
	WorkerChan    chan *worker.Worker
	StopChan      chan interface{}
	FreeWorkerBuf *worker.WorkerPool
	AllWorkerBuf  *worker.WorkerPool
	TaskBuf       *worker.TaskQueue
}

// ScheduleLoop runs the scheduling algorithm inside a goroutine
func (s *Scheduler) ScheduleLoop() {
	log.Printf("Scheduler: Scheduling loop starts\n")
loop:
	for {
		select {
		case newTask := <-s.TaskChan:
			// Receive a new task and do scheduling
			// To be implemented
		case w := <-s.WorkerChan:
			// A worker becomes free
			// To be implemented
		case <-s.StopChan:
			// Receive signal to stop scheduling
			// To be implemented
		}
	}
	log.Printf("Scheduler: Task processor ends\n")
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	go s.ScheduleLoop()
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.StopChan <- 0
	for _, w := range s.AllWorkerBuf.Pool {
		w.StopChan <- 0
	}
}
