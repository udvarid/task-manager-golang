package model

type MyTask struct {
	ID    int    `json:"id"`
	Task  string `json:"task"`
	Owner string `json:"owner"`
}
