package model

import (
	"time"
)

type TaskStatus string
type Priority string

const (
	TaskStatusPending    TaskStatus = "TODO"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
)

const (
	PriorityLow    Priority = "LOW"
	PriorityMedium Priority = "MEDIUM"
	PriorityHigh   Priority = "HIGH"
)

type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Status      TaskStatus `json:"status"`
	Priority    Priority   `json:"priority"`
	User        *User      `json:"user"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	Deleted     bool       `json:"deleted"`
}

type TaskFilter struct {
	Status   *TaskStatus `json:"status,omitempty"`
	Search   *string     `json:"search,omitempty"`
	Priority *Priority   `json:"priority,omitempty"`
}

type CreateTaskInput struct {
	Title       string      `json:"title"`
	Description *string     `json:"description,omitempty"`
	Status      *TaskStatus `json:"status"`
	Priority    *Priority   `json:"priority"`
}

type UpdateTaskInput struct {
	Title       *string     `json:"title,omitempty"`
	Description *string     `json:"description,omitempty"`
	Status      *TaskStatus `json:"status,omitempty"`
	Priority    *Priority   `json:"priority,omitempty"`
}
