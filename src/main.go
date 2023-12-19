package main

import (
	"blocking-queue/blocking_queue"
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	TASK_A_T = 5 * time.Second
	TASK_B_T = 10 * time.Second
	TASK_C_T = 15 * time.Second
	TASK_D_T = 20 * time.Second

	PROCESSOR_COUNT = 4

	SHUTDOWN_SIGNAL = 'X'
)

type processor struct {
	id      int
	queue   *blocking_queue.BlockingQueue // Use your BlockingQueue type from blocking_queue package
	sched_t int64
	work_t  int64
	real_t  int64
	wait_t  int64
}

func taskA() {
	fmt.Println("Task A starting...")
	time.Sleep(TASK_A_T)
	fmt.Println("Task A ending...")
}

func taskB() {
	fmt.Println("Task B starting...")
	time.Sleep(TASK_B_T)
	fmt.Println("Task B ending...")
}

func taskC() {
	fmt.Println("Task C starting...")
	time.Sleep(TASK_C_T)
	fmt.Println("Task C ending...")
}

func taskD() {
	fmt.Println("Task D starting...")
	time.Sleep(TASK_D_T)
	fmt.Println("Task D ending...")
}

func processorRun(p *processor, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		t, err := p.queue.Get()
		if err != nil || t.Typ == SHUTDOWN_SIGNAL {
			return
		}

		realStart := time.Now()
		switch t.Typ {
		case 'A':
			taskA()
		case 'B':
			taskB()
		case 'C':
			taskC()
		case 'D':
			taskD()
		}
		realEnd := time.Now()

		p.wait_t += realStart.Sub(t.StartTime).Milliseconds()
		p.real_t += realEnd.Sub(realStart).Milliseconds()
		p.sched_t -= realEnd.Sub(realStart).Milliseconds()

		switch t.Typ {
		case 'A':
			p.sched_t += TASK_A_T.Milliseconds()
			p.work_t += TASK_A_T.Milliseconds()
		case 'B':
			p.sched_t += TASK_B_T.Milliseconds()
			p.work_t += TASK_B_T.Milliseconds()
		case 'C':
			p.sched_t += TASK_C_T.Milliseconds()
			p.work_t += TASK_C_T.Milliseconds()
		case 'D':
			p.sched_t += TASK_D_T.Milliseconds()
			p.work_t += TASK_D_T.Milliseconds()
		}
	}
}

func main() {
	tasksAndTimes := "ABCD0123456789" // default arg
	if len(os.Args) >= 2 {
		// Use the provided command-line argument as tasksAndTimes
		tasksAndTimes = os.Args[1]
	}

	startTime := time.Now()

	if len(os.Args) < 2 {
		fmt.Println("Missing / Wrong arguments.")
		os.Exit(1)
	}

	var processors [PROCESSOR_COUNT]processor
	var wg sync.WaitGroup

	for i := 0; i < PROCESSOR_COUNT; i++ {
		processors[i].id = i
		processors[i].queue = blocking_queue.NewBlockingQueue() // Initialize your blocking queue here
		wg.Add(1)
		go processorRun(&processors[i], &wg)
	}

	for i := 0; i < len(tasksAndTimes); i++ {
		taskType := tasksAndTimes[i]

		switch taskType {
		case 'A', 'B', 'C', 'D':
			t := &blocking_queue.Task{
				Typ:       taskType,
				StartTime: time.Now(),
				EndTime:   time.Time{},
			}
			processors[i%PROCESSOR_COUNT].queue.Put(t)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			seconds := int(taskType - '0')
			time.Sleep(time.Duration(seconds) * time.Second)
		}
	}

	for i := 0; i < PROCESSOR_COUNT; i++ {
		poisonPillTask := &blocking_queue.Task{
			Typ:       SHUTDOWN_SIGNAL,
			StartTime: time.Time{},
			EndTime:   time.Time{},
		}
		processors[i].queue.Put(poisonPillTask)
	}

	wg.Wait()

	for i := 0; i < PROCESSOR_COUNT; i++ {
		p := &processors[i]
		fmt.Printf("Processor %d: Real T: %d Work T: %d Wait T: %d\n",
			p.id, p.real_t, p.work_t, p.wait_t)
	}

	fmt.Println("Elapsed:", time.Since(startTime).Seconds())
}
