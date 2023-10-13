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
	task_type string
	name string
	min_time int
	max_time int

}
func NewTask (a []string) *Task{
	if len(a) != 5{
			fmt.Println("Ошибка! Неверные аргументы функции, ожидается 3 аргумента, получено", len(a)-2)
			return nil
	}
	mi,err := strconv.Atoi(a[3])
	if err != nil {
		fmt.Println("Ошибка! Неверные аргументы функции, 2 аргумент должен быть числом")
		return nil
	}
	ma,err := strconv.Atoi(a[4])
	if err != nil {
		fmt.Println("Ошибка! Неверные аргументы функции, 3 аргумент должен быть числом")
		return nil
	}
	if mi>ma{
		fmt.Println("Ошибка! Неверные аргументы функции, минимальное время выполнения функции превосходит максимальное")
		return nil
	}
	task := Task{task_type: a[0],name: a[2], min_time: mi, max_time: ma }
	return &task
}
type Tasks struct { // переделать остальное
    secondly []*Task
    minutely []*Task
    hourly    []*Task
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
			nextTask := NewTask(strings.Split(s.Text()," "))
			if nextTask != nil{
				if nextTask.task_type == "secondly"{
					tmp_tasks.secondly = append(tmp_tasks.secondly, nextTask)
				}else if nextTask.task_type == "minutely"{
					tmp_tasks.minutely = append(tmp_tasks.minutely, nextTask)
				}else if nextTask.task_type == "hourly"{
					tmp_tasks.hourly = append(tmp_tasks.hourly,nextTask)
				}
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


func TaskMaker(task *Tasks, mu *sync.Mutex, data_channel chan *Task){
// функция передачи задач на выполнение с определенной переодичностью
	tickSecond := time.NewTicker(time.Second)
	tickMinute := time.NewTicker(2 * time.Hour)//временое значение, чтоб запиуить тикер с 00 секунд
	tickHour := time.NewTicker(2 * time.Hour)
	now := time.Now()
    nextMinutelyTick := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()+1, 0, 0, now.Location())
    nextHourlyTick := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())

    timeToNextMinutelyTick := time.Until(nextMinutelyTick)
    timeToNextHourlyTick := time.Until(nextHourlyTick)
	go func(){
		time.Sleep(timeToNextMinutelyTick)
		tickMinute = time.NewTicker(time.Minute)
	}()
	go func(){
		time.Sleep(timeToNextHourlyTick)
		tickHour = time.NewTicker(time.Hour)
	}()
	
	defer tickSecond.Stop() // освободим ресурсы, при завершении работы функции
	defer tickMinute.Stop()
	defer tickHour.Stop()
	for {
		select {
		case <-tickSecond.C:
			mu.Lock()
			for _,task :=range (task.secondly){
				data_channel<-task
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
func worker(i int, data_channel chan *Task, wg *sync.WaitGroup ){
	defer wg.Done()
	for task := range data_channel{
		// не понял как следить за закрытием канала
		a := []string{task.name,strconv.Itoa(task.min_time),strconv.Itoa(task.max_time)}
		testtasks.Wait(i,a)
	}
	
}
func TaskDoing(numberOfSimultaneousRequests int, data_channel chan *Task, wg *sync.WaitGroup ){
	for i := 0; i < numberOfSimultaneousRequests; i++{
		wg.Add(1)
		go worker(i+1, data_channel, wg)
	}
}



func main() {
	var tasks Tasks
	mu := new(sync.Mutex) 
	wg := new(sync.WaitGroup)
	numberOfSimultaneousRequests := 5 //константа определенная условиями задания
	data_channel := make(chan *Task)
	defer close(data_channel)
	go ReadData("tasks.txt", &tasks, mu)

	go TaskMaker(&tasks, mu, data_channel)
	
	TaskDoing(numberOfSimultaneousRequests, data_channel, wg)
	
	wg.Wait()

}
