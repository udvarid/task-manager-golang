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

type Configuration struct {
	Mail_psw      string `json:"mail_psw"`
	Mail_from     string `json:"mail_from"`
	FlyIo         string `json:"flyIo"`
	Environment   string `json:"environment"`
	RemoteAddress string `json:"remote_address"`
}
