package confparser

import (
	"DanilNaum/task_scheduler/internal/task"
	"bufio"
	"fmt"
	"log"
	"os"
	// "strconv"
	"strings"
	"sync"
	"time"
)

func Parse(fileName string, tasks *task.Tasks, mu *sync.Mutex) {
	//функция обновления данных из файла tasks.txt,
	// данные записываются в 3 массива в зависимости
	// от частоты выполнения и обновляется при изменении файла, проверка происходит каждые 5 секунд
	ticker := time.NewTicker(5 * time.Second)
	updateData := func() {
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatalf("Не удалость открыть файл")
		}
		defer file.Close()
		s := bufio.NewScanner(file)
		var tmp_tasks task.Tasks

		for s.Scan() {
			nextTask := task.NewTask(strings.Split(s.Text(), " "))
			if nextTask != nil {
				if nextTask.Task_type == "secondly" {
					tmp_tasks.Secondly = append(tmp_tasks.Secondly, nextTask)
				} else if nextTask.Task_type == "minutely" {
					tmp_tasks.Minutely = append(tmp_tasks.Minutely, nextTask)
				} else if nextTask.Task_type == "hourly" {
					tmp_tasks.Hourly = append(tmp_tasks.Hourly, nextTask)
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

	checkFileChanges := func() {
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