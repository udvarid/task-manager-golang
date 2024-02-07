package service

import (
	"github.com/udvarid/task-manager-golang/model"
)

var taskList = []model.MyTask{
	model.MyTask{ID: 1, Task: "Kaja készítés"},
	model.MyTask{ID: 2, Task: "Takarítés"},
	model.MyTask{ID: 3, Task: "Kutya sétáltatás"},
}

func GetAllTasks() []model.MyTask {
	return taskList
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
