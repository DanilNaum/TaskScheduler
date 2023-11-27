package task

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Task struct {
	TaskType string
	Name     string
	MinTime  int
	MaxTime  int
}
type Tasks struct {
	Secondly []*Task
	Minutely []*Task
	Hourly   []*Task
}

func NewTask(a []string) *Task {
	if len(a) != 5 {
		fmt.Println("Ошибка! Неверная запись в кофигурационном файле, ожидается 5 аргументов через пробел, получено", len(a))
		return nil
	}
	mi, err := strconv.Atoi(a[3])
	if err != nil {
		fmt.Println("Ошибка! Неверные аргументы функции, 2 аргумент должен быть числом")
		return nil
	}
	ma, err := strconv.Atoi(a[4])
	if err != nil {
		fmt.Println("Ошибка! Неверные аргументы функции, 3 аргумент должен быть числом")
		return nil
	}
	if mi > ma {
		fmt.Println("Ошибка! Неверные аргументы функции, минимальное время выполнения функции превосходит максимальное")
		return nil
	}
	
	return &Task{
		TaskType: a[0],
		Name: a[2],
		MinTime: mi,
		MaxTime: ma,
	}
}

func TransferToChan(tasks []*Task, dataChannel chan *Task, mu *sync.Mutex) {
	mu.Lock()
	for _, task := range tasks {
		dataChannel <- task
	}
	mu.Unlock()
}

func Manage(task *Tasks, mu *sync.Mutex, dataChannel chan *Task, stopAppChan chan struct{}) {
	// функция передачи задач на выполнение с определенной переодичностью
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-stopAppChan:
			stopAppChan <- struct{}{}
			return
		case <-ticker.C:
			TransferToChan(task.Secondly, dataChannel, mu)
			curTime := time.Now()
			if curTime.Second() == 0 {
				TransferToChan(task.Minutely, dataChannel, mu)
				if curTime.Minute() == 0 {
					TransferToChan(task.Hourly, dataChannel, mu)
				}
			}
		}
	}
}
