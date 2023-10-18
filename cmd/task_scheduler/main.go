package main

import (
	"DanilNaum/task_scheduler/internal/confparser"
	"DanilNaum/task_scheduler/internal/workers"
	"DanilNaum/task_scheduler/internal/task"
	"sync"
)

func main() {
	var tasks task.Tasks
	mu := new(sync.Mutex) 
	wg := new(sync.WaitGroup)
	numberOfSimultaneousRequests := 5 //константа определенная условиями задания
	data_channel := make(chan *task.Task)
	defer close(data_channel)

	go confparser.Parse("configtasks.txt", &tasks, mu)
	
	go task.Manage(&tasks, mu, data_channel)
	
	workers.Run(numberOfSimultaneousRequests, data_channel, wg)

	wg.Wait()
}
