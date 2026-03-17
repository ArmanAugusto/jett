package main

import "time"

type Priority string

const (
	PriorityHigh   Priority = "High"
	PriorityMedium Priority = "Medium"
	PriorityLow    Priority = "Low"
)

type Status string

const (
	StatusPending Status = "Pending"
	StatusDone    Status = "Done"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	StartDate time.Time `json:"start_date"`
	DueDate   time.Time `json:"due_date"`
	Priority  Priority  `json:"priority"`
	Status    Status    `json:"status"`
}

type TaskList struct {
	Tasks []Task `json:"tasks"`
}

