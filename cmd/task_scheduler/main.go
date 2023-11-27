package main

import (
	"DanilNaum/task_scheduler/internal/confparser"
	"DanilNaum/task_scheduler/internal/task"
	"DanilNaum/task_scheduler/internal/workers"
	"sync"
	
)

func main() {
	var tasks task.Tasks
	mu := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	numberOfSimultaneousRequests := 5 //константа определенная условиями задания
	stopAppChan := make(chan struct{}, 1)
	defer close(stopAppChan)
	dataChannel := make(chan *task.Task)
	defer close(dataChannel)

	go confparser.Parse("configtasks.txt", &tasks, mu, stopAppChan)

	go task.Manage(&tasks, mu, dataChannel, stopAppChan)

	workers.Run(numberOfSimultaneousRequests, dataChannel, wg,stopAppChan)
	
	wg.Wait()
	
}
