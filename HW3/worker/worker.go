package worker

import (
	"log"
	"sync"
	"time"
	c "CSE237B/HW3/constant"
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
	FinishChan chan *Worker
	PreemptChan chan interface{}
	CurrentTask *task.Task
}

// TaskProcessLoop processes tasks without preemption
func (w *Worker) TaskProcessLoop() {
	log.Printf("Worker<%d>: Task processor starts\n", w.WorkerID)
loop:
	for {
		select {
		case t := <-w.TaskChan:
			// This worker receives a new task to run
			w.CurrentTask = t
			if c.EN_PREEMPT {
				w.ProcessPreempt(t)
			} else {
				w.Process(t)
			}
			w.FinishChan <- w
		case <-w.StopChan:
			// Receive signal to stop
			break loop
		}
	}
	log.Printf("Worker<%d>: Task processor ends\n", w.WorkerID)
}

// Process runs a task on a worker without preemption
func (w *Worker) Process(t *task.Task) {
	log.Printf("Worker <%d>: App<%s>/Task<%d> starts (ddl %v)\n", w.WorkerID, t.AppID, t.TaskID, t.Deadline)
	// Process the task
	time.Sleep(t.TotalRunTime)
	// To be implemented
	log.Printf("Worker <%d>: App<%s>/Task<%d> ends\n", w.WorkerID, t.AppID, t.TaskID)
}

// Process runs a task on a worker with preemption
func (w *Worker) ProcessPreempt(t *task.Task) {
	if t.RunTime == 0 {
		log.Printf("Worker <%d>: App<%s>/Task<%d> starts (ddl %v)\n", w.WorkerID, t.AppID, t.TaskID, t.Deadline)
	} else {
		log.Printf("Worker <%d>: App<%s>/Task<%d> resumes\n", w.WorkerID, t.AppID, t.TaskID)
	}
loop:
	for {
		select {
		case <-w.PreemptChan:
			// preempted
			log.Printf("Worker <%d>: App<%s>/Task<%d> is preempted\n", w.WorkerID, t.AppID, t.TaskID)
			break loop
		default:
			time.Sleep(c.CHECK_PREEMPT_INTERVAL)
			t.RunTime += c.CHECK_PREEMPT_INTERVAL
			if t.RunTime >= t.TotalRunTime {
				// Task is done
				log.Printf("Worker <%d>: App<%s>/Task<%d> ends\n", w.WorkerID, t.AppID, t.TaskID)
				break loop
			}
		}		
	}
}
