package main

import (
	"DanilNaum/task_scheduler/internal/testtasks"
	// "DanilNaum/task_scheduler/internal/task"
	"bufio"
	"fmt"
	"strconv"

	// "fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)
type Task struct{
	name string
	min_time int
	max_time int

}
func NewTask (a []string) *Task{
	mi,_ := strconv.Atoi(a[1])
	ma,_ := strconv.Atoi(a[2])
	task := Task{name: a[0], min_time: mi, max_time: ma }
	return &task
}
type Tasks struct {
    secondly []*Task
    minutely []string
    hourly    []string
}

func ReadData(fileName string, tasks *Tasks, mu *sync.Mutex){
//функция обновления данных из файла tasks.txt,
// данные записываются в 3 массива в зависимости 
// от частоты выполнения и обновляется при изменении файла, проверка происходит каждые 30 секунд  
	ticker := time.NewTicker(30 * time.Second)
	updateData := func(){
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatalf("Не удалость открыть файл")
		}
		defer file.Close()
		s := bufio.NewScanner(file)
		var tmp_tasks Tasks

		for s.Scan(){
			if strings.Split(s.Text()," ")[0] == "secondly"{
				
				// tmp_tasks.secondly = append(tmp_tasks.secondly,  s.Text())
				nextTask := NewTask(strings.Split(s.Text()," ")[2:])
				tmp_tasks.secondly = append(tmp_tasks.secondly, nextTask)
			
			}else if strings.Split(s.Text()," ")[0] == "minutely"{
				// minute_task_tmp = append(minute_task_tmp, s.Text())
				tmp_tasks.minutely = append(tmp_tasks.minutely, s.Text())
			}else if strings.Split(s.Text()," ")[0] == "hourly"{
				tmp_tasks.hourly = append(tmp_tasks.hourly, s.Text())
			}
		}
		mu.Lock()
		*tasks = tmp_tasks
		mu.Unlock()
		
	}
	fileInfo, err := os.Stat(fileName)
    if err != nil {
        log.Fatalf("Не удалость узнать информацию о файле")
    }
	lastModifiedTime := fileInfo.ModTime()
	updateData()

	checkFileChanges := func(){
		fileInfo, err := os.Stat(fileName)
    	if err != nil {
        	log.Fatalf("Не удалость узнать информацию о файле")
    	}
		if fileInfo.ModTime().After(lastModifiedTime) {
			fmt.Println(" Файл изменен, данные обновлены")
			updateData()
			lastModifiedTime = fileInfo.ModTime()
		}
	}
	for range ticker.C {
    	checkFileChanges()
	}
}


func TaskMaker(task *Tasks, mu *sync.Mutex,data_channel chan string, data_channel2 chan *Task){
// функция передачи задач на выполнение с определенной переодичностью
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
			for _,task :=range (task.secondly){
				data_channel2<-task
			}
			mu.Unlock()
		case <-tickMinute.C:
			mu.Lock()
			for _,task :=range (task.minutely){
				data_channel<-task
			}
			mu.Unlock()
		case <-tickHour.C:
			mu.Lock()
			for _,task :=range (task.hourly){
				data_channel<-task
			}
			mu.Unlock()
		
		}
	}
}

func TaskDoing(numberOfSimultaneousRequests int, data_channel chan string, data_channel2 chan *Task  ){
	thread := make(map[int] chan struct{})
	for i := 0; i < numberOfSimultaneousRequests; i++ {
		thread[i] = make(chan struct{},1)
		thread[i] <-  struct{}{}
		defer close(thread[i])
	}

	for task := range data_channel2{
		// a := strings.Split(s," ")
		a := []string{task.name,strconv.Itoa(task.min_time),strconv.Itoa(task.max_time)}
		go func(a []string){			
			select {
			case <-thread[0]:
				testtasks.Wait(1,a)
				thread[0] <-  struct{}{}
			case <-thread[1]:
				testtasks.Wait(2,a)
				thread[1] <-  struct{}{}	
			case <-thread[2]:
				testtasks.Wait(3,a)
				thread[2] <-  struct{}{}
			case <-thread[3]:
				testtasks.Wait(4,a)
				thread[3] <-  struct{}{}	
			case <-thread[4]:
				testtasks.Wait(5,a)
				thread[4] <-  struct{}{}	
			}		
		}(a)
	}

}

func main() {
	// var second_task []string
	// var minute_task []string
	// var hour_task []string
	var tasks Tasks
	mu := new(sync.Mutex) 
	
	numberOfSimultaneousRequests := 5 //константа определенная условиями задания
	data_channel := make(chan string) // канал используемый для передачи очереди задач
	data_channel2 := make(chan *Task)
	defer close(data_channel)
	defer close(data_channel2)
	// go ReadData("tasks.txt", &tasks ,&second_task, &minute_task, &hour_task, mu)
	go ReadData("tasks.txt", &tasks, mu)

	// go TaskMaker(&second_task, &minute_task, &hour_task, mu, data_channel)
	go TaskMaker(&tasks, mu, data_channel, data_channel2)
	
	TaskDoing(numberOfSimultaneousRequests, data_channel, data_channel2)
	// Наличие структуры на thread означает, что поток свободен и можно передать на него следующую задачу
	

}
