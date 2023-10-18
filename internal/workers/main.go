package workers

import (
	// "DanilNaum/task_scheduler/internal/confparser"
	"DanilNaum/task_scheduler/internal/task"
	"DanilNaum/task_scheduler/internal/testtasks"
	"strconv"
	"sync"
)

func worker(i int, data_channel chan *task.Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range data_channel {
		// не понял как следить за закрытием канала
		a := []string{task.Name, strconv.Itoa(task.Min_time), strconv.Itoa(task.Max_time)}
		testtasks.Wait(i, a)
	}

}
func Run(numberOfSimultaneousRequests int, data_channel chan *task.Task, wg *sync.WaitGroup) {
	for i := 0; i < numberOfSimultaneousRequests; i++ {
		wg.Add(1)
		go worker(i+1, data_channel, wg)
	}
}