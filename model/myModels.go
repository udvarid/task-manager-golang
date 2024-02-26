package model

import "time"

type MyTask struct {
	ID    int    `json:"id"`
	Task  string `json:"task"`
	Owner string `json:"owner"`
}

type SessionWithTime struct {
	Session   string
	SessDate  time.Time
	IsChecked bool
}
