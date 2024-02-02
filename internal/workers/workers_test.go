package workers_test

import (
	"DanilNaum/task_scheduler/internal/task"
	"DanilNaum/task_scheduler/internal/workers"
	"sync"
	"testing"
	"time"
)

func TestRun(t *testing.T){
	testTable := []struct{
		numberOfSimultaneousRequests int
		dataChannel chan *task.Task
		wg *sync.WaitGroup
		stopAppChan chan struct{}
		exaTask []*task.Task
	}{
		{
			numberOfSimultaneousRequests: 2,
			dataChannel:make(chan *task.Task),
			wg: new(sync.WaitGroup),
			stopAppChan: make(chan struct{}, 1),
			exaTask: []*task.Task{
				{TaskType: "secondly", Name: "secondly1",MinTime: 1,MaxTime: 2},
				{TaskType: "secondly", Name: "secondly2",MinTime: 6,MaxTime: 7},
			},
		},
	}
	
	for _, testCase := range(testTable) {
		defer close(testCase.dataChannel)
		go func ()  {
			for _,task := range(testCase.exaTask){
			
			   testCase.dataChannel<-task
		   }
		}()
		
		
		go func(stopAppChan chan struct{}){
			
			time.Sleep(5 * time.Second)
			close(stopAppChan)
		}(testCase.stopAppChan)
		go func(){
			
			time.Sleep(8 * time.Second)

			t.Errorf("Процесс не завершен вовремя")
			
		}()


		workers.Run(testCase.numberOfSimultaneousRequests, testCase.dataChannel,testCase.wg,testCase.stopAppChan)
		// time.Sleep(5 * time.Second)
		// close(testCase.stopAppChan)
		testCase.wg.Wait()
	}
}