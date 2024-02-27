package service

import (
	"sort"
	"time"

	"github.com/udvarid/task-manager-golang/model"
	"github.com/udvarid/task-manager-golang/repository/taskRepository"
)

func GetAllTasks(owner string) ([]model.MyTaskDto, []model.MyTaskDto) {
	tasks := taskRepository.GetAllTask(owner)
	var taskDtos []model.MyTaskDto
	timeNow := time.Now()
	for _, task := range tasks {
		overDue := false
		if timeNow.After(task.Deadline) {
			overDue = true
		}
		taskDto := model.MyTaskDto{
			ID:       task.ID,
			Task:     task.Task,
			Deadline: task.Deadline,
			DeadLStr: task.Deadline.Format(time.DateOnly),
			Warning:  overDue,
			Owner:    task.Owner,
		}
		taskDtos = append(taskDtos, taskDto)
	}
	sort.Slice(taskDtos, func(i, j int) bool {
		return taskDtos[i].Deadline.Before(taskDtos[j].Deadline)
	})
	var taskDtosOverDue []model.MyTaskDto
	var taskDtosNormal []model.MyTaskDto
	for _, taskDto := range taskDtos {
		if taskDto.Warning {
			taskDtosOverDue = append(taskDtosOverDue, taskDto)
		} else {
			taskDtosNormal = append(taskDtosNormal, taskDto)
		}
	}
	return taskDtosOverDue, taskDtosNormal
}

func DeleteTask(taskId int) {
	taskRepository.DeleteTask(taskId)
}

func AddTask(task model.NewTask, owner string) {
	taskRepository.AddTask(task, owner)
}
