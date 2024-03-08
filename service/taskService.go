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
			ID:             task.ID,
			Task:           task.Task,
			Deadline:       task.Deadline,
			DeadLStr:       task.Deadline.Format(time.DateOnly),
			Warning:        overDue,
			Owner:          task.Owner,
			LastWarnedTime: task.LastWarnedTime,
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

func DeleteTask(taskId int, owner string) {
	myTask := taskRepository.GetTask(taskId)
	if myTask.Owner == owner {
		taskRepository.DeleteTask(taskId)
	}

}

func ProlongTask(taskId int, prolongDays int, owner string) {
	if prolongDays > 0 {
		myTask := taskRepository.GetTask(taskId)
		if myTask.Owner == owner {
			prolongedDeadLine := myTask.Deadline.AddDate(0, 0, prolongDays)
			myTask.Deadline = prolongedDeadLine
			taskRepository.UpdateTask(taskId, myTask)
		}
	}
}

func AddTask(task model.NewTask, owner string) {
	taskRepository.AddTask(task, owner)
}

func UpdateWarningDate(id int) {
	task := taskRepository.GetTask(id)
	task.LastWarnedTime = time.Now()
	taskRepository.UpdateTask(id, task)
}
