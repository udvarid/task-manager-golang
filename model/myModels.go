package model

import "time"

type MyTask struct {
	ID       int       `json:"id"`
	Task     string    `json:"task"`
	Deadline time.Time `json:"deadLine"`
	Owner    string    `json:"owner"`
}

type MyTaskDto struct {
	ID       int       `json:"id"`
	Task     string    `json:"task"`
	Deadline time.Time `json:"deadLine"`
	DeadLStr string    `json:"deadLStr"`
	Warning  bool      `json:"warning"`
	Owner    string    `json:"owner"`
}

type SessionWithTime struct {
	Session   string
	SessDate  time.Time
	IsChecked bool
}

type NewTask struct {
	Task     string `json:"task"`
	Deadline string `json:"deadline"`
}
