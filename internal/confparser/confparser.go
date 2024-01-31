package confparser

import (
	"DanilNaum/task_scheduler/internal/task"
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func Parse(fileName string, tasks *task.Tasks, mu *sync.Mutex, stopAppChan chan struct{}) {
	//функция обновления данных из файла tasks.txt,
	// данные записываются в 3 массива в зависимости
	// от частоты выполнения и обновляется при изменении файла, проверка происходит каждые 5 секунд
	ticker := time.NewTicker(5 * time.Second)

	// fileInfo, err := os.Stat(fileName)
	// if err != nil {
	// 	log.Fatalf("Не удалость узнать информацию о файле")
	// }
	// lastModifiedTime := fileInfo.ModTime()
	lastModifiedTime := time.Unix(0, 0)
	updateData(fileName, tasks, mu,stopAppChan)
	for {
		fileInfo, err := os.Stat(fileName)
		if err != nil {
			fmt.Print("Не удалость узнать информацию о файле")
			// stopAppChan <- struct{}{}
			close(stopAppChan)
			return
		}
		if fileInfo.ModTime().After(lastModifiedTime) {
			fmt.Println(" Файл изменен, данные обновлены")
			updateData(fileName, tasks, mu,stopAppChan)
			lastModifiedTime = fileInfo.ModTime()
		}
		select {
			case <-stopAppChan:
				return
			case <-ticker.C:
				continue
		}
	}

}

func updateData(fileName string, tasks *task.Tasks, mu *sync.Mutex,stopAppChan chan struct{}) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Print("Не удалость открыть файл")
		// stopAppChan <- struct{}{}
		close(stopAppChan)
		return
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	var tmpTasks task.Tasks

	for s.Scan() {
		nextTask := task.NewTask(strings.Split(s.Text(), " "))
		if nextTask != nil {
			if nextTask.TaskType == "secondly" {
				tmpTasks.Secondly = append(tmpTasks.Secondly, nextTask)
			} else if nextTask.TaskType == "minutely" {
				tmpTasks.Minutely = append(tmpTasks.Minutely, nextTask)
			} else if nextTask.TaskType == "hourly" {
				tmpTasks.Hourly = append(tmpTasks.Hourly, nextTask)
			}
		}
	}
	mu.Lock()
	*tasks = tmpTasks
	mu.Unlock()
}
