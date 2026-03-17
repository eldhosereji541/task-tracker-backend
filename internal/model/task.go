package model

import (
	"time"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "TODO"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
)

type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	Deleted     bool       `json:"deleted"`
}

type TaskFilter struct {
	Status *TaskStatus `json:"status,omitempty"`
	Search *string     `json:"search,omitempty"`
}
