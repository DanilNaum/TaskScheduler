package task_test

import (
	"DanilNaum/task_scheduler/internal/task"
	
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestManage(t *testing.T) {
	// Arrange
	tasks := &task.Tasks{Secondly: []*task.Task{
		{TaskType: "secondly", Name: "1st s task",MinTime: 5,MaxTime: 6},
		{TaskType: "secondly", Name: "2st s task",MinTime: 1,MaxTime: 2},
	},Minutely: []*task.Task{
		{TaskType: "minutely", Name: "1st m task",MinTime: 5,MaxTime: 6},
		{TaskType: "minutely", Name: "2st m task",MinTime: 1,MaxTime: 2},
	},Hourly: []*task.Task{
		{TaskType: "hourly", Name: "1st h task",MinTime: 5,MaxTime: 6},
	},}
	mu := new(sync.Mutex)
	stopAppChan := make(chan struct{}, 1)
	dataChannel := make(chan *task.Task)
	defer close(dataChannel)
	// Act

	go task.Manage(tasks,mu,dataChannel,stopAppChan)

	// Assert
	timer := time.After(60*time.Second)
	prevSecTime := time.Now()
	for {
		select{
		case <-timer:
			// stopAppChan <- struct{}{}
			close(stopAppChan)
			return
		case task, ok := <-dataChannel:
			if !ok {
				t.Errorf("Канал datachan сломан ")
				return
			}
			if strings.Contains(task.TaskType,  "secondly"){
				timeLeft := time.Since(prevSecTime).Milliseconds()
				if !((900 <timeLeft && timeLeft < 1100) || (timeLeft < 100)){
					t.Errorf("Задача передана в канал в некорректное вреямя, ошибка в передаче секунд, "+ strconv.Itoa(int(timeLeft))+" милесекунд с прошлого запроса")
					return
				}
				prevSecTime = time.Now()
			}else if strings.Contains(task.TaskType,  "minutely"){
				timeNow := time.Now().Second()
				if  !(timeNow == 0){
					t.Errorf("Задача передана в канал в некорректное вреямя, ошибка в передаче минут,  сейчас"  + strconv.Itoa(int(timeNow))+" секунд, запуск ожидается в начале каждой минуты ")
					return
				}
			}else if strings.Contains(task.TaskType,  "hourly"){
				timeNowMin := time.Now().Minute()
				timeNowSec := time.Now().Second()
				
				if  !(timeNowMin + timeNowSec == 0){
					t.Errorf("Задача передана в канал в некорректное вреямя, ошибка в передаче часов,  сейчас"  + strconv.Itoa(int(timeNowMin)) +":" +strconv.Itoa(int(timeNowSec))+" минут:секунд, запуск ожидается в начале каждого час")
					return
				}
			}

		}

	}
	//... тест должен проверять, что таски попадают своевременно в канал
}
func TestNewTask(t *testing.T) {
	t.Run("must_return_nil_on_incorrect_arguments_number", func(t *testing.T) {
		testTable := []struct{
			a []string
			res *task.Task
		}{
			{
				a: []string{"secondly", "testtasks.Wait", "secondly1", "1", "2"},
				res: &task.Task{
					TaskType: "secondly",
					Name:  "secondly1",
					MinTime: 1,
					MaxTime: 2,},
			},
			{
				a: []string{"testtasks.Wait", "secondly1", "1", "2"},
				res: nil,
			},
			{
				a: []string{ "secondly1", "1", "2"},
				res: nil,
			},
			{
				a: []string{ "1", "2"},
				res: nil,
			},
			{
				a: []string{ "2"},
				res: nil,
			},
		}
		for _,tastCase:=range(testTable){
			res := task.NewTask(tastCase.a)
			if res != nil && tastCase.res!=nil{
				if *res!=*tastCase.res{
					t.Errorf("Функция NewTask вернула не то значение, которое ожидалось. Ожидалось %+v, получено %+v",*res,*tastCase.res)
				}else if (res!=nil&&tastCase.res==nil)||(res==nil&&tastCase.res!=nil){
						t.Errorf("Функция NewTask вернула не то значение, которое ожидалось,  одно из значений ожидаемое или полученое равняется nil, а другое нет" )
				}
			}
			
		}
		
		// ... тест проверяет, что логируется ошибка, если длина слайса не равна 5
	})
	t.Run("must_return_nil_if_4rd_or_5th_arguments_is_not_number", func(t *testing.T) {
		testTable := []struct{
			a []string
			res *task.Task
		}{
			{
				a: []string{"secondly", "testtasks.Wait", "secondly1", "1", "2"},
				res: &task.Task{
					TaskType: "secondly",
					Name:  "secondly1",
					MinTime: 1,
					MaxTime: 2,},
			},
			{
				a: []string{"secondly", "testtasks.Wait", "secondly1", "abob", "2"},
				res: nil,
			},
			{
				a: []string{"secondly", "testtasks.Wait", "secondly1", "1", "vfv5"},
				
				res: nil,
			},
			{
				a: []string{"secondly", "testtasks.Wait", "secondly1", "dfg1", "dfg2"},
				res: nil,
			},
		}
		for _,tastCase:=range(testTable){
			res := task.NewTask(tastCase.a)
			if res != nil && tastCase.res!=nil{
				if *res!=*tastCase.res{
					t.Errorf("Функция NewTask вернула не то значение, которое ожидалось. Ожидалось %+v, получено %+v",*res,*tastCase.res)
				}else if (res!=nil&&tastCase.res==nil)||(res==nil&&tastCase.res!=nil){
						t.Errorf("Функция NewTask вернула не то значение, которое ожидалось,  одно из значений ожидаемое или полученое равняется nil, а другое нет" )
				}
			}
			
		}
		
		// ... тест проверяет, что логируется ошибка, если 4 или 5 аргумент не числа
	})
	t.Run("must_return_nil_if_5th_arguments_less_than_4rd_is_not_number", func(t *testing.T) {
		testTable := []struct{
			a []string
			res *task.Task
		}{
			{
				a: []string{"secondly", "testtasks.Wait", "secondly1", "1", "2"},
				res: &task.Task{
					TaskType: "secondly",
					Name:  "secondly1",
					MinTime: 1,
					MaxTime: 2,},
			},
			{
				a: []string{"secondly", "testtasks.Wait", "secondly1", "6", "2"},
				res: nil,
			},
			
		}
		for _,tastCase:=range(testTable){
			res := task.NewTask(tastCase.a)
			if res != nil && tastCase.res!=nil{
				if *res!=*tastCase.res{
					t.Errorf("Функция NewTask вернула не то значение, которое ожидалось. Ожидалось %+v, получено %+v",*res,*tastCase.res)
				}else if (res!=nil&&tastCase.res==nil)||(res==nil&&tastCase.res!=nil){
						t.Errorf("Функция NewTask вернула не то значение, которое ожидалось,  одно из значений ожидаемое или полученое равняется nil, а другое нет" )
				}
			}
			
		}
		
		// ... тест проверяет, что логируется ошибка, если 4 аргумент больше 5 
	})
}
