package main

import (
	"DanilNaum/task_scheduler/internal/testtasks"
	"bufio"
	"time"
	"os"
	"strings"
	"sync"
)

func ReadData(fileName string, second_task *[]string, minute_task *[]string, hour_task *[]string, mu *sync.Mutex){//функция обновления данных из файла tasks.txt, данные записываются в 3 массива в зависимости от частоты выполнения и обновляются при каждой новой итерации
	//возможно это не очень оптимально и нужно менять данные только при изменении файла. 
	//Так же может быть допустима задержка, тогда обновлять файл можно с некоторой заранее заданой переодичностью
	
	for{
		dir, err := os.Getwd() //не разобрался как указывать относительный путь к файлу, поэтому пришлось запрашивать дирректорию
		if err != nil {
			panic(err)
		}
		file, err := os.Open(dir + "/cmd/task_scheduler/" + fileName)
		if err != nil {
			panic(err)
		}
		s := bufio.NewScanner(file)
		var second_task_tmp []string
		var minute_task_tmp []string
		var hour_task_tmp []string
		for s.Scan(){
			if strings.Split(s.Text()," ")[0] == "secondly"{
				second_task_tmp = append(second_task_tmp, s.Text())
			}else if strings.Split(s.Text()," ")[0] == "minutely"{
				minute_task_tmp = append(minute_task_tmp, s.Text())
			}else if strings.Split(s.Text()," ")[0] == "hourly"{
				hour_task_tmp = append(hour_task_tmp, s.Text())
			}
		}
		mu.Lock()
		*second_task = second_task_tmp
		*minute_task = minute_task_tmp
		*hour_task = hour_task_tmp
		mu.Unlock()
		file.Close()
	}
}

func TaskMaker(second_task *[]string, minute_task *[]string, hour_task *[]string, mu *sync.Mutex,data_channel chan string){ 
	// функция передачи задач на выполнение с определенной переодичностью
	defer close(data_channel)
	tickSecond := time.NewTicker(time.Second)
	tickMinute := time.NewTicker(time.Minute)
	tickHour := time.NewTicker(time.Hour)
	defer tickSecond.Stop() // освободим ресурсы, при завершении работы функции
	defer tickMinute.Stop()
	defer tickHour.Stop()
	for {
		select {
		case <-tickSecond.C:
			mu.Lock()
			for _,task :=range (*second_task){
				data_channel<-task
			}
			mu.Unlock()
		case <-tickMinute.C:
			mu.Lock()
			for _,task :=range (*minute_task){
				data_channel<-task
			}
			mu.Unlock()
		case <-tickHour.C:
			mu.Lock()
			for _,task :=range (*hour_task){
				data_channel<-task
			}
			mu.Unlock()
		
		}
	}
}

func TaskDoing(numberOfSimultaneousRequests int, data_channel chan string  ){
	thread := make(map[int] chan struct{})
	for i := 0; i < numberOfSimultaneousRequests; i++ {
		thread[i] = make(chan struct{},1)
		thread[i] <-  struct{}{}
		defer close(thread[i])
	}

	for s := range data_channel{
		a := strings.Split(s," ")
		//wg.Add(1)
		// go func(wg *sync.WaitGroup, a []string){
		go func(a []string){
			// defer wg.Done()
			
			select {
			case <-thread[0]:
				testtasks.Wait(1,a[2:])
				thread[0] <-  struct{}{}
			case <-thread[1]:
				testtasks.Wait(2,a[2:])
				thread[1] <-  struct{}{}	
			case <-thread[2]:
				testtasks.Wait(3,a[2:])
				thread[2] <-  struct{}{}
			case <-thread[3]:
				testtasks.Wait(4,a[2:])
				thread[3] <-  struct{}{}	
			case <-thread[4]:
				testtasks.Wait(5,a[2:])
				thread[4] <-  struct{}{}	
			}		
		// }(wg,a)
		}(a)
	}

}

func main() {
	var second_task []string
	var minute_task []string
	var hour_task []string
	mu := new(sync.Mutex) 
	//wg := new(sync.WaitGroup)
	numberOfSimultaneousRequests := 5 //константа определенная условиями задания
	data_channel := make(chan string) // канал используемый для передачи очереди задач

	go ReadData("tasks.txt", &second_task, &minute_task, &hour_task, mu)

	// go func() { 
		
	// }()
	go TaskMaker(&second_task, &minute_task, &hour_task, mu, data_channel)
	
	TaskDoing(numberOfSimultaneousRequests, data_channel)
	// Наличие структуры на thread означает, что поток свободен и можно передать на него следующую задачу
	
	//wg.Wait()
}
