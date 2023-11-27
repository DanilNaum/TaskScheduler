package workers

import (
	// "DanilNaum/task_scheduler/internal/confparser"
	"DanilNaum/task_scheduler/internal/task"
	"DanilNaum/task_scheduler/internal/testtasks"
	"fmt"

	// "strconv"
	"sync"
)

func worker(i int, dataChannel chan *task.Task, wg *sync.WaitGroup, stopAppChan chan struct{}) {
	defer func() {
		fmt.Printf("Воркер %d завершился\n", i)
		wg.Done()
	}()

	for {
		select {
		case <-stopAppChan:
			stopAppChan <- struct{}{}
			return
		case task, ok := <-dataChannel:
			if !ok {
				return
			}
			testtasks.Wait(i, task)

		}

	}

}
func Run(numberOfSimultaneousRequests int, dataChannel chan *task.Task, wg *sync.WaitGroup, stopAppChan chan struct{}) {

	for i := 0; i < numberOfSimultaneousRequests; i++ {
		wg.Add(1)
		go worker(i+1, dataChannel, wg, stopAppChan)
	}
}
