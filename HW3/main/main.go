package main

import (
	c "CSE237B/HW3/constant"
	"fmt"
	"CSE237B/HW3/scheduler"
	"CSE237B/HW3/task"
	"time"
)

func main() {
	apps := []*task.App{}

	// Task specifications
	taskSpecs := []task.TaskSpec{
		task.TaskSpec{
			Period:           4 * time.Second,
			TotalRunTimeMean: 1 * time.Second,
			TotalRunTimeStd:  300 * time.Millisecond,
			RelativeDeadline: 3 * time.Second,
		},
		task.TaskSpec{
			Period:           2 * time.Second,
			TotalRunTimeMean: 1 * time.Second,
			TotalRunTimeStd:  300 * time.Millisecond,
			RelativeDeadline: 2 * time.Second,
		},
	}

	// Create and initialize the scheduler
	sched := scheduler.NewScheduler()

	// Create all applications
	for i, taskSpec := range taskSpecs {
		apps = append(apps, task.NewApp(fmt.Sprintf("app%d", i), taskSpec))
	}

	
	// To be implemented, initialization process

	// Start the scheduler
	sched.Start()

	// Start all applications
	for _, app := range apps {
		app.TaskChan = sched.TaskChan
		app.Start()
	}

	time.Sleep(c.TEST_TIME)

	// Stop all applications
	for _, app := range apps {
		app.Stop()
	}

	// Stop the scheduler
	sched.Stop()
}
