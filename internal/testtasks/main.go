package testtasks

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)
func Wait(thread int,args []string){ // использован лишний аргумент thread, который не предпологался в задаче, не придумал как передать поток на котором выполняется функция другим способом
	if len(args) == 3{
		min_time,err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Ошибка! Неверные аргументы функции, 2 аргумент должен быть числом")
				return
			}
		max_time,err := strconv.Atoi(args[2])
			if err != nil {
				fmt.Println("Ошибка! Неверные аргументы функции, 3 аргумент должен быть числом")
				return
			}
		if min_time < max_time{
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "Thread #",thread,": starting job\"",args[0],"\"")
			dalay_time := rand.Intn( max_time-min_time+1) + min_time
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "Thread #",thread,": wating for", dalay_time,"seconds...")
			time.Sleep(time.Duration(dalay_time)*time.Second)
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "Thread #",thread,":job is complited")
		} else {
			fmt.Println("Ошибка! Неверные аргументы функции, минимальное время выполнения функции превосходит максимальное")
			// Обработка передачи неверных аргументов функции, выбор способа реакции зависит от цели, предпологаю, что программа должна продолжить выполнение
			// panic("Неверные аргументы функции..")
			return
		}
	} else{
		fmt.Println("Ошибка! Неверные аргументы функции, ожидается 3 аргумента, получено", len(args))
		// Обработка передачи неверных аргументов функции, выбор способа реакции зависит от цели, предпологаю, что программа должна продолжить выполнение
		// panic("Неверные аргументы функции..")
		return

	}
}