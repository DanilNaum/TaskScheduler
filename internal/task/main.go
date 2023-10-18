package task

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)
type Task struct{
	Task_type string
	Name string
	Min_time int
	Max_time int

}
type Tasks struct { 
    Secondly []*Task
    Minutely []*Task
    Hourly    []*Task
}
func NewTask (a []string) *Task{
	if len(a) != 5{
			fmt.Println("Ошибка! Неверная запись в кофигурационном файле, ожидается 5 аргументов через пробел, получено", len(a))
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
	task := Task{Task_type: a[0],Name: a[2], Min_time: mi, Max_time: ma }
	return &task
}

func TransferToChan(task []*Task,data_channel chan *Task, mu *sync.Mutex){
	mu.Lock()
	for _, task := range task {
		data_channel <- task
	}
	mu.Unlock()
}

func Manage(task *Tasks, mu *sync.Mutex, data_channel chan *Task) {
	// функция передачи задач на выполнение с определенной переодичностью
	tickSecond := time.NewTicker(time.Second)
	tickMinute := time.NewTicker(2 * time.Hour) //временое значение, чтоб запиуить тикер с 00 секунд
	tickHour := time.NewTicker(2 * time.Hour)
	now := time.Now()
	nextMinutelyTick := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()+1, 0, 0, now.Location())
	nextHourlyTick := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())

	timeToNextMinutelyTick := time.Until(nextMinutelyTick)
	timeToNextHourlyTick := time.Until(nextHourlyTick)
	go func() {
		time.Sleep(timeToNextMinutelyTick)
		tickMinute = time.NewTicker(time.Minute)
		TransferToChan(task.Minutely,data_channel,mu)
	}()
	go func() {
		time.Sleep(timeToNextHourlyTick)
		tickHour = time.NewTicker(time.Hour)
		TransferToChan(task.Hourly,data_channel,mu)
	}()

	defer tickSecond.Stop() // освободим ресурсы, при завершении работы функции
	defer tickMinute.Stop()
	defer tickHour.Stop()
	for {
		select {
		case <-tickSecond.C:
			TransferToChan(task.Secondly,data_channel,mu)
		case <-tickMinute.C:
			TransferToChan(task.Minutely,data_channel,mu)
		case <-tickHour.C:
			TransferToChan(task.Hourly,data_channel,mu)

		}
	}
}