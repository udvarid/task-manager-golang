package repository

import (
	"github.com/udvarid/task-manager-golang/model"
)

var taskList = []model.MyTask{
	{ID: 1, Task: "Kaja készítés", Owner: "donat1977"},
	{ID: 2, Task: "Takarítás", Owner: "donat1977"},
	{ID: 3, Task: "Kutya sétáltatás", Owner: "donat1977"},
}

func GetAllTask(owner string) []model.MyTask {
	var result []model.MyTask

	for _, task := range taskList {
		if task.Owner == owner {
			result = append(result, task)
		}
	}

	return result
}

func DeleteTask(taskId int) {
	if index := getIndexOfTask(taskId); index != -1 {
		taskList = append(taskList[:index], taskList[index+1:]...)
	}
}

func getIndexOfTask(taskId int) int {
	for i, t := range taskList {
		if t.ID == taskId {
			return i
		}
	}
	return -1
}

func AddTask(task string, owner string) {
	newTask := model.MyTask{
		ID:    findNextId(),
		Task:  task,
		Owner: owner,
	}
	taskList = append(taskList, newTask)
}

func findNextId() int {
	maxId := 0
	for _, task := range taskList {
		if maxId < task.ID {
			maxId = task.ID
		}
	}
	return maxId + 1
}
