package confparser_test

import (
	"DanilNaum/task_scheduler/internal/confparser"
	"DanilNaum/task_scheduler/internal/task"
	"fmt"
	// "reflect"
	"sync"
	"testing"
	"time"
)
func TestParse(t *testing.T){
	type data struct {
		fileName string
		tasks *task.Tasks
		mu *sync.Mutex
		stopAppChan chan struct{}
		res *task.Tasks
		stopTestingChan chan struct{}
		
	}
	testTable :=[]data{
		{
			fileName:"E:/golang_progect/TaskScheduler/internal/confparser/test1.txt", //файла не существует
			tasks : &task.Tasks{Secondly: []*task.Task{},
				Minutely: []*task.Task{},
				Hourly: []*task.Task{},},
			mu : new(sync.Mutex),
			stopAppChan : make(chan struct{}, 1),
			res: nil,
			stopTestingChan : make(chan struct{}, 1),
		},
		{
			fileName:"E:/golang_progect/TaskScheduler/internal/confparser/test.txt", //файла существует
			tasks : &task.Tasks{Secondly: []*task.Task{},
				Minutely: []*task.Task{},
				Hourly: []*task.Task{},},
			mu : new(sync.Mutex),
			stopAppChan : make(chan struct{}, 1),
			res: &task.Tasks{Secondly: []*task.Task{
				{TaskType: "secondly", Name: "secondly1",MinTime: 1,MaxTime: 2},
			},Minutely:  []*task.Task{},
			Hourly: []*task.Task{},
			},
			stopTestingChan : make(chan struct{}, 1),
		},
	}
	
	for i,testCase := range(testTable){
		fmt.Printf("\ntest_case %d\n",i)
		timer := time.After(30*time.Second)
		
		// stopTestingChan := make(chan struct{}, 1)
		ticker := time.NewTicker(10 * time.Second)
		
		go confparser.Parse(testCase.fileName,testCase.tasks,testCase.mu,testCase.stopAppChan)
		go func(testCase data){
			for{
				select{
				case <-testCase.stopAppChan:
					if testCase.res != nil{
						t.Errorf("Канал закрыт, хотя не должен")
						
					}
					close(testCase.stopTestingChan)
					return
				case <- timer:
					close(testCase.stopAppChan)
					close(testCase.stopTestingChan)
					return
				}
			
			}
		}(testCase)
		func(testCase data){
		for {
			select{
			case <-testCase.stopTestingChan:
				return
			case <-ticker.C:
				if testCase.res != nil &&  testCase.tasks!= nil{
					if len(testCase.tasks.Secondly) ==len(testCase.res.Secondly){
						for i := range(testCase.tasks.Secondly){
							if *testCase.tasks.Secondly[i] != *testCase.res.Secondly[i]{
								t.Errorf("Считаные значения не соответствуют содержанию файла. Ожидалось %+v, получено %+v", *testCase.tasks.Secondly[i] ,*testCase.res.Secondly[i])
								close(testCase.stopTestingChan)
							}
						}
					}else{
						t.Errorf("Считаные значения не соответствуют содержанию файла." )
						close(testCase.stopTestingChan)
					}
					if len(testCase.tasks.Minutely) ==len(testCase.res.Minutely){
						for i := range(testCase.tasks.Minutely){
							if *testCase.tasks.Minutely[i] != *testCase.res.Minutely[i]{
								t.Errorf("Считаные значения не соответствуют содержанию файла. Ожидалось %+v, получено %+v", *testCase.tasks.Minutely[i] ,*testCase.res.Minutely[i])
								close(testCase.stopTestingChan)
							}
						}
					}else{
						t.Errorf("Считаные значения не соответствуют содержанию файла." )
						close(testCase.stopTestingChan)
					}
					if len(testCase.tasks.Hourly) ==len(testCase.res.Hourly){
						for i := range(testCase.tasks.Hourly){
							if *testCase.tasks.Hourly[i] != *testCase.res.Hourly[i]{
								t.Errorf("Считаные значения не соответствуют содержанию файла. Ожидалось %+v, получено %+v", *testCase.tasks.Hourly[i] ,*testCase.res.Hourly[i])
								close(testCase.stopTestingChan)
							}
						}
					}else{
						t.Errorf("Считаные значения не соответствуют содержанию файла." )
						close(testCase.stopTestingChan)
					}
				}
			}

		}
	}(testCase)
		
	}
	
}