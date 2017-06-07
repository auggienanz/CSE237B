package scheduler

import (
	"log"
	"sort"
	"sync"
	c "CSE237B/HW3/constant"
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

type ByDeadline []*task.Task

func (queue ByDeadline) Len() int {
	return len(queue)
}

func (queue ByDeadline) Swap(i, j int) {
	queue[i], queue[j] = queue[j], queue[i]
}

func (queue ByDeadline) Less(i, j int) bool {
	return queue[j].Deadline.After(queue[i].Deadline)
}

// ScheduleLoop runs the scheduling algorithm inside a goroutine
func (s *Scheduler) ScheduleLoop() {
	log.Printf("Scheduler: Scheduling loop starts\n")
loop:
	for {
		select {
		case newTask := <-s.TaskChan:
			log.Printf("Scheduler: Received new task\n")
			// Add task to task buffer
			s.TaskBuf.Lock.Lock()
			s.TaskBuf.Queue = append(s.TaskBuf.Queue, newTask)
			sort.Sort(ByDeadline(s.TaskBuf.Queue))
			// If there's a free worker, assign this task
			s.FreeWorkerBuf.Lock.Lock()
			if len(s.FreeWorkerBuf.Pool) != 0 {
				s.FreeWorkerBuf.Pool[0].TaskChan <- s.TaskBuf.Queue[0]
				s.FreeWorkerBuf.Pool = s.FreeWorkerBuf.Pool[1:]
				s.TaskBuf.Queue = s.TaskBuf.Queue[1:]
			} else if c.EN_PREEMPT {
				// check all currently running tasks
				for _,w := range(s.AllWorkerBuf.Pool) {
					// If the current task has a further deadline, preempt
					if (s.TaskBuf.Queue[0].Deadline.Before(w.CurrentTask.Deadline)) {
						w.PreemptChan <- 0
					}
				}
			}
			s.FreeWorkerBuf.Lock.Unlock()
			s.TaskBuf.Lock.Unlock()
		case w := <-s.WorkerChan:
			// A worker becomes free
			s.TaskBuf.Lock.Lock()
			
			// Check if its task finished. If not, we need to add the task back to the queue
			if (w.CurrentTask.TotalRunTime > w.CurrentTask.RunTime) {
				s.TaskBuf.Queue = append(s.TaskBuf.Queue, w.CurrentTask)
				sort.Sort(ByDeadline(s.TaskBuf.Queue))
			}
			// If there's a task in the queue, assign it to this worker
			if (len(s.TaskBuf.Queue) != 0) {
				w.TaskChan <- s.TaskBuf.Queue[0]
				s.TaskBuf.Queue = s.TaskBuf.Queue[1:]
			} else {
				// Otherwise, add this worker to the free worker queue
				s.FreeWorkerBuf.Lock.Lock()
				s.FreeWorkerBuf.Pool = append(s.FreeWorkerBuf.Pool, w)
				s.FreeWorkerBuf.Lock.Unlock()
			}
			s.TaskBuf.Lock.Unlock()
			


		case <-s.StopChan:
			// Receive signal to stop scheduling
			// Wait for all workers to finish tasks
			s.FreeWorkerBuf.Lock.Lock()
			s.AllWorkerBuf.Lock.Lock()
			for (len(s.FreeWorkerBuf.Pool) != len(s.AllWorkerBuf.Pool)) {
				w := <-s.WorkerChan
				s.FreeWorkerBuf.Pool = append(s.FreeWorkerBuf.Pool, w)
			}
			s.FreeWorkerBuf.Lock.Unlock()
			s.AllWorkerBuf.Lock.Unlock()
			break loop
		}
	}
	log.Printf("Scheduler: Task processor ends\n")
}

func NewScheduler() *Scheduler{
	var s Scheduler
	s.FreeWorkerBuf = &worker.WorkerPool{Pool: make([]*worker.Worker,0), Lock: new(sync.Mutex)}
	s.AllWorkerBuf = &worker.WorkerPool{Pool: make([]*worker.Worker,0), Lock: new(sync.Mutex)}
	s.TaskBuf = &worker.TaskQueue{Queue: make([]*task.Task,0), Lock: new(sync.Mutex)}
	s.TaskChan = make(chan *task.Task)
	s.WorkerChan = make(chan *worker.Worker)
	s.StopChan = make(chan interface{})
	// Create workers
	s.FreeWorkerBuf.Lock.Lock()
	s.AllWorkerBuf.Lock.Lock()
	var i int
	for i = 0; i < c.WORKER_NR; i++ {
		w := &worker.Worker{ 
			WorkerID: i,
			TaskChan: make(chan *task.Task),
			StopChan: make(chan interface{}),
			FinishChan: s.WorkerChan,
			PreemptChan: make(chan interface{}),
		}
		s.AllWorkerBuf.Pool = append(s.AllWorkerBuf.Pool, w)
		s.FreeWorkerBuf.Pool = append(s.FreeWorkerBuf.Pool, w)
		go w.TaskProcessLoop()
	}
	s.FreeWorkerBuf.Lock.Unlock()
	s.AllWorkerBuf.Lock.Unlock()
	return &s
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
