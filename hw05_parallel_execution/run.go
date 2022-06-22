package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksToGo := uint32(len(tasks) - 1)
	availableGoroutines := uint32(n)
	maxErrors := uint32(m)

	var wg sync.WaitGroup
	var errorCount uint32 = 0
	i := 0

	for tasksToGo > 0 {
		if maxErrors > 0 && errorCount >= maxErrors {
			return ErrErrorsLimitExceeded
		}

		if availableGoroutines > 0 {
			if i > len(tasks)-1 {
				break
			}
			nextTask := tasks[i]
			i++
			atomic.AddUint32(&availableGoroutines, ^uint32(0))
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := nextTask()
				if err != nil {
					atomic.AddUint32(&errorCount, 1)
				}
				atomic.AddUint32(&tasksToGo, ^uint32(0))
				atomic.AddUint32(&availableGoroutines, 1)
			}()
		}
	}

	wg.Wait()
	if maxErrors > 0 && errorCount >= maxErrors {
		return ErrErrorsLimitExceeded
	}

	return nil
}
