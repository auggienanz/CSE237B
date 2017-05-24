package worker

import (
	"log"
	"sync"
	"CSE237B/HW3/task"
)

type TaskQueue struct {
	Queue []*task.Task
	Lock  *sync.Mutex
}

type WorkerPool struct {
	Pool []*Worker
	Lock *sync.Mutex
}

// Worker is the agent to process tasks
type Worker struct {
	WorkerID int
	TaskChan chan *task.Task
	StopChan chan interface{}
}

// TaskProcessLoop processes tasks without preemption
func (w *Worker) TaskProcessLoop() {
	log.Printf("Worker<%d>: Task processor starts\n", w.WorkerID)
loop:
	for {
		select {
		case t := <-w.TaskChan:
			// This worker receives a new task to run
			// To be implemented
		case <-w.StopChan:
			// Receive signal to stop
			// To be implemented
			break loop
		}
	}
	log.Printf("Worker<%d>: Task processor ends\n", w.WorkerID)
}

// Process runs a task on a worker without preemption
func (w *Worker) Process(t *task.Task) {
	log.Printf("Worker <%d>: App<%s>/Task<%d> starts (ddl %v)\n", w.WorkerID, t.AppID, t.TaskID, t.Deadline)
	// Process the task
	// To be implemented
	log.Printf("Worker <%d>: App<%s>/Task<%d> ends\n", w.WorkerID, t.AppID, t.TaskID)
}
