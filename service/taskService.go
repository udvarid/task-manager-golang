package service

import (
	"github.com/udvarid/task-manager-golang/model"
	"github.com/udvarid/task-manager-golang/repository/taskRepository"
)

func GetAllTasks(owner string) []model.MyTask {
	return taskRepository.GetAllTask(owner)
}

func DeleteTask(taskId int) {
	taskRepository.DeleteTask(taskId)
}

func AddTask(task string, owner string) {
	taskRepository.AddTask(task, owner)
}
