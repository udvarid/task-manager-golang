package service

import (
	"github.com/udvarid/task-manager-golang/model"
	"github.com/udvarid/task-manager-golang/repository"
)

func GetAllTasks(owner string) []model.MyTask {
	return repository.GetAllTask(owner)
}

func DeleteTask(taskId int) {
	repository.DeleteTask(taskId)
}

func AddTask(task string, owner string) {
	repository.AddTask(task, owner)
}
