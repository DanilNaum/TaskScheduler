package testtasks

import (
	"DanilNaum/task_scheduler/internal/task"
	"fmt"
	"math/rand"
	"time"
)

func Wait(thread int, task *task.Task) {
	minTime := task.MinTime
	maxTime := task.MaxTime
	name := task.Name
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "Thread #", thread, ": starting job\"", name, "\"")
	dalayTime := rand.Intn(maxTime-minTime+1) + minTime
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "Thread #", thread, ": wating for", dalayTime, "seconds...")
	time.Sleep(time.Duration(dalayTime) * time.Second)
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "Thread #", thread, ":job is complited")
}
